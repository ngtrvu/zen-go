package zen

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
)

type TestClient struct {
	Writer        httptest.ResponseRecorder
	RouterContext *chi.Context
	Context       context.Context
}

func NewTestClient(ctx context.Context) *TestClient {
	return &TestClient{
		Writer:        *httptest.NewRecorder(),
		RouterContext: chi.NewRouteContext(),
		Context:       ctx,
	}
}

func (t *TestClient) MakeRequest(method string, path string, payload io.Reader) *http.Request {
	req := httptest.NewRequest(method, path, payload)
	req.Header.Set("Content-Type", "application/json")

	req = req.WithContext(t.Context)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, t.RouterContext))

	return req
}

func (t *TestClient) ResponseJSON() (Response, error) {
	var resp Response
	err := json.Unmarshal(t.Writer.Body.Bytes(), &resp)
	if err != nil {
		return resp, errors.New("failed to parse response")
	}

	return resp, nil
}

func (t *TestClient) ResetRecorder() {
	t.Writer = *httptest.NewRecorder()
	t.RouterContext = chi.NewRouteContext()
}
