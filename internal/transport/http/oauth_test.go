package httpapi

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tetafro/nott-backend-go/internal/auth"
)

func TestOAuthController(t *testing.T) {
	t.Run("Get list of providers", func(t *testing.T) {
		c := NewOAuthController(
			map[string]*auth.OAuthProvider{
				"example-1": {Name: "example-1", URL: "http://example-1.com"},
				"example-2": {Name: "example-2", URL: "http://example-2.com"},
			},
			nil, nil, nil,
		)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)

		c.Providers(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusOK)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.JSONEq(t, string(body), `{
			"data": [
				{
					"name": "example-1",
					"url": "http://example-1.com"
				},
				{
					"name": "example-2",
					"url": "http://example-2.com"
				}
			]
		}`)
	})
}
