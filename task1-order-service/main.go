package main

import (
	"log"
	"net/http"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"example.com/order-service/httptransport"
	"example.com/order-service/implementation"
	"example.com/order-service/repository"
	"example.com/order-service/transport"
)

func main() {
	logger := kitlog.NewLogfmtLogger(os.Stderr)
	logger = kitlog.NewSyncLogger(logger)
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
	logger = level.NewFilter(logger, level.AllowInfo())

	repo := repository.NewMemoryRepository()
	svc := implementation.NewService(repo, logger)
	handler := httptransport.NewService(transport.MakeEndpoints(svc), logger)
	level.Info(logger).Log("msg", "order service started", "addr", ":8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
