package httpapi

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespond(t *testing.T) {
	t.Run("respond with data", func(t *testing.T) {
		w := httptest.NewRecorder()

		respond(w, http.StatusOK, map[string]string{"field": "value"})

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusOK)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.JSONEq(t, string(body), `{
			"data": {
				"field": "value"
			}
		}`)
	})

	t.Run("respond with error", func(t *testing.T) {
		w := httptest.NewRecorder()

		respond(w, http.StatusBadRequest, "bad request")

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.JSONEq(t, string(body), `{"error": "bad request"}`)
	})

	t.Run("shortcut: bad request", func(t *testing.T) {
		w := httptest.NewRecorder()

		badRequest(w, "bad request")

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.JSONEq(t, string(body), `{"error": "bad request"}`)
	})

	t.Run("shortcut: not found", func(t *testing.T) {
		w := httptest.NewRecorder()

		notFound(w)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusNotFound)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.JSONEq(t, string(body), `{"error": "not found"}`)
	})

	t.Run("shortcut: unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()

		unauthorized(w)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusUnauthorized)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.JSONEq(t, string(body), `{"error": "unauthorized"}`)
	})

	t.Run("shortcut: internal server error", func(t *testing.T) {
		w := httptest.NewRecorder()

		internalServerError(w)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.JSONEq(t, string(body), `{"error": "internal server error"}`)
	})
}
