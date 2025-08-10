package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"status": http.StatusOK,
		"data":   "Ok!",
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Something Went Wrong!")
		return
	}
}
