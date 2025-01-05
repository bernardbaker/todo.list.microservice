package domain

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const secretKey = "your_shared_secret_key"

func validateHMAC(method, path, timestamp, body, receivedSig string) error {
	message := fmt.Sprintf("%s\n%s\n%s\n%s", method, path, timestamp, body)

	// Compute HMAC
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))
	expectedSig := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Compare signatures
	if !hmac.Equal([]byte(receivedSig), []byte(expectedSig)) {
		return errors.New("invalid HMAC signature")
	}

	// Validate timestamp (5-minute window)
	reqTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return errors.New("invalid timestamp format")
	}
	if time.Since(reqTime) > 20*time.Second {
		return errors.New("request timestamp is outside the allowed window")
	}

	return nil
}

func HmacInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Extract metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	// Get required fields
	method := info.FullMethod
	fmt.Println(method)
	timestamp := md.Get("x-timestamp")[0]
	receivedSig := md.Get("authorization")[0]

	// Convert req to JSON string
	reqMessage, ok := req.(proto.Message)
	if !ok {
		return nil, errors.New("request does not implement proto.Message")
	}

	// Use protojson.MarshalOptions for compact JSON
	marshalOptions := protojson.MarshalOptions{
		Indent:          "",   // No indentation
		EmitUnpopulated: true, // Include default values for unset fields (optional)
	}
	reqJSON, err := marshalOptions.Marshal(reqMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request to JSON: %v", err)
	}

	body := string(reqJSON) // Extract request body if needed

	// Validate HMAC
	if err := validateHMAC("POST", method, timestamp, body, receivedSig); err != nil {
		return nil, err
	}

	// Proceed with the request
	return handler(ctx, req)
}
