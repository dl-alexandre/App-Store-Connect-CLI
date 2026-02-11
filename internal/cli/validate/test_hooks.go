package validate

import (
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

// SetClientFactory replaces the ASC client factory for tests.
// It returns a restore function to reset the previous handler.
func SetClientFactory(fn func() (*asc.Client, error)) func() {
	previous := clientFactory
	if fn == nil {
		clientFactory = shared.GetASCClient
	} else {
		clientFactory = fn
	}
	return func() {
		clientFactory = previous
	}
}
