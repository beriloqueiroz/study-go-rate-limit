package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	config "github.com/beriloqueiroz/study-go-rate-limit/configs"
	routes "github.com/beriloqueiroz/study-go-rate-limit/internal/infra/web/routes/api"
	webserver "github.com/beriloqueiroz/study-go-rate-limit/internal/infra/web/server"
)

func main() {

	// graceful exit
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	initCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// load environment configs
	configs, err := config.LoadConfig([]string{"."})
	if err != nil {
		panic(err)
	}

	server := webserver.NewWebServer(configs.WebServerPort, 2)
	route := routes.NewTestSimpleRoute()
	server.AddRoute("GET /", route.Handler)
	srvErr := make(chan error, 1)
	go func() {
		fmt.Println("Starting web server on port", configs.WebServerPort)
		srvErr <- server.Start()
	}()

	// Wait for interruption.
	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-initCtx.Done():
		log.Println("Shutting down due to other reason...")
	}
}

func rateLimit(key string, ip string) {
	if key != "" {
		fmt.Println("Consider key")
		return
	}
	fmt.Println("Consider ip")
}
