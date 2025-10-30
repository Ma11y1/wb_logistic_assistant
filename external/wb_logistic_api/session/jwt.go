package session

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func decodeJWT(jwt string, target interface{}) error {
	parts := strings.Split(jwt, ".")
	if len(parts) < 2 {
		return fmt.Errorf("invalid token format")
	}

	payload := parts[1]
	paddedPayload := payload + strings.Repeat("=", (4-len(payload)%4)%4)
	decodedPayload, err := base64.URLEncoding.DecodeString(paddedPayload)
	if err != nil {
		return fmt.Errorf("failed to decode payload: %v", err)
	}

	if err = json.Unmarshal(decodedPayload, &target); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	return nil
}
