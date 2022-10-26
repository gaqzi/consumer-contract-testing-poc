package pablo_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {

	// Create Pact connecting to local Daemon
	pact := &dsl.Pact{
		Provider: "MyProvider",
	}
	pactDir := "../pablosdk/pacts/"

	// Start provider API in the background
	httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))

	// Verify the Provider using the locally saved Pact Files
	_, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: "http://localhost:8000",
		PactURLs:        []string{filepath.ToSlash(fmt.Sprintf("%s/pablosdk-pablo.json", pactDir))},
		StateHandlers: types.StateHandlers{
			// Setup any state required by the test
			// in this case, we ensure there is a "user" in the system
			"A valid test authorization token": func() error {
				return nil
			},
		},
	})

	require.NoError(t, err, "failed to verify provider")
}
