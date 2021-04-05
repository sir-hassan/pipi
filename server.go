package main

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const apiPort = 8080

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	level.Info(logger).Log("msg", "starting pipi...")

	router := mux.NewRouter()
	router.HandleFunc("/movie/amazon/{amazon_id}", createHandler(logger))

	http.Handle("/", router)
	level.Info(logger).Log("msg", "listening started", "port", apiPort)
	ln, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(apiPort))
	if err != nil {
		level.Error(logger).Log("msg", "listening on port failed", "port", apiPort, "error", err)
		return
	}
	server := &http.Server{}

	go func() {
		err := server.Serve(ln)
		if err != http.ErrServerClosed {
			level.Error(logger).Log("msg", "closing server failed", "error", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	level.Info(logger).Log("msg", "signal received", "sig", sig)
	level.Info(logger).Log("msg", "terminating server...")

	if err := server.Shutdown(context.Background()); err != nil {
		level.Error(logger).Log("msg", "terminating server failed", "error", err)
	} else {
		level.Info(logger).Log("msg", "terminating server succeeded")
	}
}

func createHandler(logger log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.WriteHeader(http.StatusOK)
		level.Debug(logger).Log("msg", "new request")
		amazonID := vars["amazon_id"]

		req, _ := http.NewRequest("GET", "https://www.amazon.de/gp/product/"+amazonID, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:86.0) Gecko/20100101 Firefox/86.0")
		httpClient := &http.Client{}
		res, err := httpClient.Do(req)
		if err != nil {
			writeReply(logger, w, 500, "internal server error")
			return
		}
		defer res.Body.Close()

	}
}

func writeReply(logger log.Logger, w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	if _, err := w.Write([]byte(message)); err != nil {
		level.Error(logger).Log("msg", "writing connection", "error", err)
	}
}
