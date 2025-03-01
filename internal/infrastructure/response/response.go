package response

import (
	"encoding/json"
	"net/http"

	rfc7807 "github.com/moogar0880/problems"
)

func Write(w http.ResponseWriter, body []byte) {
	w.Header().Add("Content-Type", "application/json")
	_, err := w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	problem := rfc7807.NewDetailedProblem(statusCode, err.Error())
	marshalledProblem, err := json.Marshal(problem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if _, err_ := w.Write(marshalledProblem); err_ != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
