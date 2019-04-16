package main

import (
	"log"
	"net"

	"github.com/enfipy/auvima/src/config"
	"github.com/enfipy/auvima/src/helpers"
	"github.com/enfipy/auvima/src/services"
)

func main() {
	cnfg := config.InitConfig()
	cnfg.Settings = helpers.GetSettings("/settings.yaml")

	if cnfg.AppEnv != "production" {
		log.SetFlags(0)
	} else {
		log.Print("Server started")
	}

	port := ":" + cnfg.ServerPort
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	echo, close := services.InitServices(cnfg)
	go helpers.GracefulShutdown(close)

	echo.Listener = lis
	if err := echo.Start(""); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
