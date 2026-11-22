package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"

	studentv1 "example.com/student-service/gen/student/v1"
	"example.com/student-service/grpctransport"
	"example.com/student-service/httptransport"
	"example.com/student-service/implementation"
	"example.com/student-service/repository"
	"example.com/student-service/transport"
)

func main() {
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	grpcAddr := flag.String("grpc.addr", ":9090", "gRPC listen address")
	flag.Parse()

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = level.NewFilter(logger, level.AllowInfo())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	repo := repository.NewMemoryRepository()
	svc := implementation.NewService(repo, logger)
	endpoints := transport.MakeEndpoints(svc)
	httpHandler := httptransport.NewService(endpoints, logger)

	httpServer := &http.Server{Addr: *httpAddr, Handler: httpHandler, ReadHeaderTimeout: 5 * time.Second}
	grpcServer := grpc.NewServer()
	studentv1.RegisterStudentServiceServer(grpcServer, grpctransport.NewService(endpoints))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	errCh := make(chan error, 2)

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr, "msg", "server started")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()
	go func() {
		listener, err := newListener(*grpcAddr)
		if err != nil {
			errCh <- err
			return
		}
		level.Info(logger).Log("transport", "gRPC", "addr", *grpcAddr, "msg", "server started")
		if err := grpcServer.Serve(listener); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		level.Info(logger).Log("msg", "shutdown signal received")
	case err := <-errCh:
		level.Error(logger).Log("msg", "server stopped unexpectedly", "err", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		level.Error(logger).Log("transport", "HTTP", "msg", "graceful shutdown failed", "err", err)
	}
	grpcServer.GracefulStop()
	level.Info(logger).Log("msg", "servers stopped")
}
