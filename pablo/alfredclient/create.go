package alfredclient

import "context"

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

func (c Implementation) Create(ctx context.Context, req CreateRequest) (CreateResponse, error) {
	return CreateResponse{}, nil
}
