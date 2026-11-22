package httptransport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"example.com/order-service/order"
	"example.com/order-service/transport"
)

func NewService(e transport.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	opts := []kithttp.ServerOption{kithttp.ServerErrorLogger(logger), kithttp.ServerErrorEncoder(encodeError)}
	r.Methods(http.MethodPost).Path("/orders").Handler(kithttp.NewServer(e.Create, decodeCreate, encodeResponse, opts...))
	r.Methods(http.MethodGet).Path("/orders/{id}").Handler(kithttp.NewServer(e.GetByID, decodeGet, encodeResponse, opts...))
	r.Methods(http.MethodPatch).Path("/orders/{id}/status").Handler(kithttp.NewServer(e.ChangeStatus, decodeChangeStatus, encodeResponse, opts...))
	return r
}

func decodeCreate(_ context.Context, r *http.Request) (interface{}, error) {
	var o order.Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		return nil, fmt.Errorf("%w: invalid JSON", order.ErrInvalidOrder)
	}
	return transport.CreateRequest{Order: o}, nil
}
func decodeGet(_ context.Context, r *http.Request) (interface{}, error) {
	return transport.GetByIDRequest{ID: mux.Vars(r)["id"]}, nil
}
func decodeChangeStatus(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("%w: invalid JSON", order.ErrInvalidOrder)
	}
	return transport.ChangeStatusRequest{ID: mux.Vars(r)["id"], Status: body.Status}, nil
}
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	if errors.Is(err, order.ErrInvalidOrder) {
		code = http.StatusBadRequest
	}
	if errors.Is(err, order.ErrOrderNotFound) {
		code = http.StatusNotFound
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
