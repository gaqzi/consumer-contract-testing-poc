package pablosdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IntentClient struct {
	pabloAddr string
	authToken string
}

type Instrument struct {
	ID          string `json:"id" pact:"example=a41x5mdfns7c9s7"`
	Description string `json:"description"`
}

type Method struct {
	Name        string       `json:"name" pact:"example=creditcard,regex=creditcard|cash|psp"`
	Instruments []Instrument `json:"instruments"`
}

type CreateRequest struct {
	Amount int64 `json:"amount" pact:"example=999"`
}

type CreateResponse struct {
	ID      string   `json:"id" pact:"example=b16872c595994147"`
	Methods []Method `json:"methods"`
}

const (
	intentURL = "/payments/intent"
)

func NewIntentClient(addr, authToken string) *IntentClient {
	return &IntentClient{pabloAddr: addr, authToken: authToken}
}

func (i *IntentClient) Create(ctx context.Context, req CreateRequest) (CreateResponse, error) {
	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(req); err != nil {
		return CreateResponse{}, fmt.Errorf("failed to encode create request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, i.pabloAddr+intentURL, buf)
	if err != nil {
		return CreateResponse{}, fmt.Errorf("failed to create create httpRequest: %w", err)
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer "+i.authToken)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return CreateResponse{}, fmt.Errorf("failed to call create intent: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != 201 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return CreateResponse{}, fmt.Errorf("failed to read create error response body: %w", err)
		}

		return CreateResponse{}, fmt.Errorf("unexpected http response '%d' from create request: %s", resp.StatusCode, body)
	}

	var cResp CreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return CreateResponse{}, fmt.Errorf("failed to decode successful http response: %w", err)
	}

	return cResp, nil
}
