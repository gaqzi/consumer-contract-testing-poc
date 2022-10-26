package alfredclient_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/dsl"

	"github.com/gaqzi/consumer-driven-contract-poc/pablo/alfredclient"
)

func TestCreate(t *testing.T) {
	// Create Pact connecting to local Daemon
	pact := &dsl.Pact{
		Consumer: "pablo",
		Provider: "alfred",
		Host:     "localhost",
	}
	defer pact.Teardown()

	// Pass in test case. This is the component that makes the external HTTP call
	var test = func() (err error) {
		c := alfredclient.New(fmt.Sprintf("http://localhost:%d", pact.Server.Port))
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := c.Auth(ctx, "valid-user"); err != nil {
			return fmt.Errorf("failed to auth: %w", err)
		}

		_, err = c.Create(ctx, alfredclient.CreateRequest{Amount: 999})

		return err
	}

	// Set up our expected interactions.
	pact.
		AddInteraction().
		Given("A valid user for auth").
		UponReceiving("A valid auth request").
		WithRequest(dsl.Request{
			Method:  "POST",
			Path:    dsl.String("/auth"),
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body: map[string]string{
				"userId": "valid-user",
			},
		}).
		WillRespondWith(dsl.Response{
			Status:  201,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body:    dsl.Match(&alfredclient.AuthResponse{}),
		})

	pact.
		AddInteraction().
		Given("A valid authToken for the authed user").
		UponReceiving("A valid create request").
		WithRequest(dsl.Request{
			Method:  "POST",
			Path:    dsl.String("/purchase/intent"),
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json"), "Authorization": dsl.String("Bearer valid-authToken-from-auth")},
			Body: map[string]int64{
				"amount": 999,
			},
		}).
		WillRespondWith(dsl.Response{
			Status:  201,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body:    dsl.Match(&alfredclient.CreateResponse{}),
		})

	// Run the test, verify it did what we expected and capture the contract
	if err := pact.Verify(test); err != nil {
		log.Fatalf("Error on Verify: %v", err)
	}
}
