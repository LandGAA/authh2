package main

import (
	"github.com/LandGAA/authh2/internal/app"
	"github.com/LandGAA/authh2/pkg/grpc/server"
	"github.com/LandGAA/authh2/pkg/logger"
)

func main() {
	logger.LoggerRun()
	go app.Run()
	go server.Run()
	select {}
}
