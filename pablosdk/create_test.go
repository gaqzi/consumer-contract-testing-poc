package pablosdk_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"

	"github.com/gaqzi/consumer-driven-contract-poc/pablosdk"
)

func TestCreate(t *testing.T) {
	// Create Pact connecting to local Daemon
	pact := &dsl.Pact{
		Consumer: "pablosdk",
		Provider: "Pablo",
		Host:     "localhost",
	}
	defer pact.Teardown()

	// Pass in test case. This is the component that makes the external HTTP call
	var test = func() (err error) {
		c := pablosdk.NewIntentClient(fmt.Sprintf("http://localhost:%d", pact.Server.Port), "valid-token")

		_, err = c.Create(context.TODO(), pablosdk.CreateRequest{Amount: 999})

		return err
	}

	// Set up our expected interactions.
	pact.
		AddInteraction().
		Given("A valid test authorization token").
		UponReceiving("A valid crate request").
		WithRequest(dsl.Request{
			Method:  "POST",
			Path:    dsl.String("/payments/intent"),
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json"), "Authorization": dsl.String("Bearer valid-token")},
			Body: map[string]int64{
				"amount": 999,
			},
		}).
		WillRespondWith(dsl.Response{
			Status:  201,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body:    dsl.Match(&pablosdk.CreateResponse{}),
		})

	// Run the test, verify it did what we expected and capture the contract
	if err := pact.Verify(test); err != nil {
		log.Fatalf("Error on Verify: %v", err)
	}
}
