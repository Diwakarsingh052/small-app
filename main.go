package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"rest-api/auth"
	"rest-api/handlers"
	"rest-api/models"
	"time"
)

func main() {
	err := startApp()
	if err != nil {
		panic(err)
	}
}

func startApp() error {

	// Initialize authentication support
	publicKey, err := SetupAuth()
	if err != nil {
		return err
	}
	a, err := auth.New(publicKey)
	if err != nil {
		return err
	}

	//----------- Initialize cache
	slog.Info("main : Started : Initializing cache ")
	c := models.NewConn()

	h, err := handlers.API(a, c)
	if err != nil {
		return err
	}
	api := http.Server{
		Addr:         ":8081",
		ReadTimeout:  500 * time.Second,
		WriteTimeout: 500 * time.Second,
		IdleTimeout:  500 * time.Second,
		Handler:      h,
	}

	serverErr := make(chan error)
	go func() {
		serverErr <- api.ListenAndServe()
	}()
	//shutdown channel intercepts ctrl+c signals
	shutdown := make(chan os.Signal, 1)
	// signal.Notify will notify the given channel when someone produces the given os signal
	signal.Notify(shutdown, os.Interrupt)

	select {
	// listening for errors that might happen during server startup, usually port is already being used
	case err := <-serverErr:
		return err
	case <-shutdown:
		fmt.Println("Gracefully shutting down server...")
		// creating a timer of 5sec for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		//Shutdown gracefully shuts down the server without interrupting any active connections.
		//Shutdown works by first closing all open listeners
		err := api.Shutdown(ctx)
		if err != nil {
			//close immediately closes all active net. Listeners and any connections in state
			// forceful close
			err := api.Close()
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func SetupAuth() (*rsa.PublicKey, error) {

	slog.Info("main : Started : Initializing authentication support")

	publicPEM, err := os.ReadFile("pubkey.pem")
	if err != nil {
		return nil, fmt.Errorf("reading auth public key %w", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return nil, fmt.Errorf("parsing auth public key %w", err)
	}
	return publicKey, nil
}
