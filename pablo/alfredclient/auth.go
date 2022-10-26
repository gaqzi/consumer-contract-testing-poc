package alfredclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AuthResponse struct {
	Token string `json:"authToken" pact:"example=mnc0ien4m2x"`
}

func (c Implementation) Auth(ctx context.Context, user string) error {
	reqBody, err := json.Marshal(map[string]string{"userId": user})
	if err != nil {
		return fmt.Errorf("failed to marshal the auth request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.addr+authURL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read auth response body: %w", err)
	}

	if resp.StatusCode != 201 {
		return fmt.Errorf("unexpected http response '%d' for auth request: %s", resp.StatusCode, body)
	}

	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return fmt.Errorf("failed to unmarshal auth response (%q): %w", body, err)
	}

	if authResp.Token == "" {
		return fmt.Errorf("received invalid (empty) token from auth response")
	}

	c.auth = authResp
	return nil
}
