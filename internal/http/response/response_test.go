package response

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type data struct {
	Text string `json:"text"`
}

func TestNew(t *testing.T) {
	r := New()

	if r.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %v, but got %v", http.StatusOK, r.StatusCode)
	}
}

func TestResponseWithStatus(t *testing.T) {
	s := http.StatusBadRequest

	r := New().WithStatus(s)

	if r.StatusCode != s {
		t.Fatalf("Expected %v, but got %v", s, r.StatusCode)
	}
}

func TestResponseWithData(t *testing.T) {
	d := data{Text: "Foo"}

	r := New().WithData(d)

	if r.Data != d {
		t.Fatalf("Expected %v, but got %v", d, r.Data)
	}
}

func TestResponseWithError(t *testing.T) {
	e := fmt.Errorf("error")

	r := New().WithError(e)

	if r.Error != e {
		t.Fatalf("Expected %v, but got %v", e, r.Error)
	}
}

func TestResponseWrite(t *testing.T) {
	d := data{Text: "Foo"}

	t.Run("Write successful response with data", func(t *testing.T) {
		r := New().WithStatus(http.StatusOK).WithData(d)

		w := httptest.NewRecorder()
		r.Write(w)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected code %v, but got %v", http.StatusOK, w.Code)
		}

		expectedBody := `{"data":{"text":"Foo"}}`
		actualBody := strings.TrimSpace(w.Body.String())
		if actualBody != expectedBody {
			t.Fatalf("Expected body %v, but got %v", expectedBody, actualBody)
		}
	})

	t.Run("Write error response", func(t *testing.T) {
		r := New().WithStatus(http.StatusBadRequest).WithError(fmt.Errorf("error"))

		w := httptest.NewRecorder()
		r.Write(w)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected code %v, but got %v", http.StatusBadRequest, w.Code)
		}

		expectedBody := `{"error":"error"}`
		actualBody := strings.TrimSpace(w.Body.String())
		if actualBody != expectedBody {
			t.Fatalf("Expected body to be %v but got %v", expectedBody, actualBody)
		}
	})
}
