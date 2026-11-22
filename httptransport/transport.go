// Package httptransport предоставляет HTTP/JSON интерфейс сервиса.
package httptransport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	kittransport "github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"

	"example.com/student-service/student"
	"example.com/student-service/transport"
)

func NewService(endpoints transport.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods(http.MethodPost).Path("/students").Handler(kithttp.NewServer(
		endpoints.Create, decodeCreateRequest, encodeResponse, options...,
	))
	r.Methods(http.MethodGet).Path("/students/{id}").Handler(kithttp.NewServer(
		endpoints.GetByID, decodeGetByIDRequest, encodeResponse, options...,
	))
	r.Methods(http.MethodPatch).Path("/students/{id}/group").Handler(kithttp.NewServer(
		endpoints.ChangeGroup, decodeChangeGroupRequest, encodeResponse, options...,
	))
	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var st student.Student
	if err := json.NewDecoder(r.Body).Decode(&st); err != nil {
		return nil, fmt.Errorf("%w: invalid JSON body", student.ErrInvalidRequest)
	}
	return transport.CreateRequest{Student: st}, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return transport.GetByIDRequest{ID: mux.Vars(r)["id"]}, nil
}

func decodeChangeGroupRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Group string `json:"group"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("%w: invalid JSON body", student.ErrInvalidRequest)
	}
	return transport.ChangeGroupRequest{ID: mux.Vars(r)["id"], Group: body.Group}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	statusCode := http.StatusInternalServerError
	switch {
	case errors.Is(err, student.ErrInvalidRequest), errors.Is(err, student.ErrInvalidStudent):
		statusCode = http.StatusBadRequest
	case errors.Is(err, student.ErrStudentNotFound):
		statusCode = http.StatusNotFound
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
