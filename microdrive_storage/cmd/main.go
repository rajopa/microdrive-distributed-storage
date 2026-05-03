package main

import (
	"log/slog"
	"microdrive_storage/internal/app"
)

func main() {
	imageServer := app.NewImageServer()
	grpcServer := app.NewGrpcServer()
	err := grpcServer.GrpcServeServer(imageServer, ":8086")
	if err != nil {
		slog.Warn("Server shutdown with error", "error", err.Error())
	}
}
