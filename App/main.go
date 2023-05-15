package main

import (
	"context"
	"errors"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const Addr = ":670"

func newRouter() *httprouter.Router {
	mux := httprouter.New()
	mux.GET("/youtube/armin", getChannelStats())
	return mux
}

func getChannelStats() httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.Write([]byte("Hello Armin !"))
	}
}

func main() {
	srv := &http.Server{
		Addr:    Addr,
		Handler: newRouter(),
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		log.Println("service interrupt received")

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("http server shut down error :%v", err)

		}
		log.Println("shutdown complete")

		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("fatal http server failed to start: %v", err)
		}

	}
	<-idleConnsClosed
	log.Println("Service Stop")
}
