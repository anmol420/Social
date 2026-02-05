package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("Internal Server Error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Bad Request Error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Not Found Error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusNotFound, "Resource Not Found")
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("Conflict Error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusConflict, "Unauthorized")
}

func (app *application) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("Unauthorized Error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusUnauthorized, err.Error())
}

func (app *application) forbiddenError(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("Forbidden Error", "method", r.Method, "path", r.URL.Path)
	writeJSONError(w, http.StatusForbidden, "Forbidden")
}

func (app *application) ratelimitExceedError(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("Rate Limit Exceeded Error", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Retry-After", retryAfter)
	writeJSONError(w, http.StatusTooManyRequests, "Rate Limit Exceeded, Try After: "+retryAfter)
}
