package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anmol420/Social/internal/auth"
	"github.com/anmol420/Social/internal/store"
	"github.com/anmol420/Social/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()
	logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockCacheStore()
	mockAuthenticator := &auth.TestAuthenticator{}
	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCacheStore,
		authenticator: mockAuthenticator,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code: %d and Acutal response code received: %d", expected, actual)
	}
}
