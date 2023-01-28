package httplisten

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// A replacement for http.ListenAndServe which catches termination signals and gracefully shuts down connections
func Serve(addr string, handler http.Handler) error {
	var srv http.Server

	return startServer(&srv, addr, handler, func() error {
		return srv.ListenAndServe()
	})
}

// A replacement for http.ListenAndServeTLS which catches termination signals and gracefully shuts down connections
func ServeTLS(addr, certFile, keyFile string, handler http.Handler) error {
	var srv http.Server

	return startServer(&srv, addr, handler, func() error {
		return srv.ListenAndServeTLS(certFile, keyFile)
	})
}

// The internal server withc graceful shutdown
func startServer(srv *http.Server, addr string, handler http.Handler, server func() error) error {
	var notifyServerShutdown chan int

	srv.Addr = addr
	if handler != nil {
		srv.Handler = handler
	}

	notifyServerShutdown = make(chan int, 1)
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()
		srv.Shutdown(ctx)
		close(notifyServerShutdown)
	}()

	err := server()
	<-notifyServerShutdown
	return err
}
