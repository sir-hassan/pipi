package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/sir-hassan/pipi/backend"
	"github.com/sir-hassan/pipi/parse"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const apiPort = 8080

// ErrorReply is used to encode error api response payload.
type ErrorReply struct {
	Error string `json:"error"`
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	level.Info(logger).Log("msg", "starting pipi...")

	pageFetcher := backend.NewWebPageFetcher(&http.Client{})
	handler := createHandler(logger, pageFetcher)

	http.Handle("/", handler)
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

func createHandler(logger log.Logger, pageFetcher backend.PageFetcher) http.Handler {
	handler := mux.NewRouter()
	handler.HandleFunc("/movie/amazon/{amazon_id}", createMovieHandleFunc(logger, pageFetcher))
	handler.NotFoundHandler = http.HandlerFunc(createNotFoundHandleFunc(logger))
	return handler
}

func createNotFoundHandleFunc(logger log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		writeErrorReply(logger, w, 404, "the requested path is not found")
	}
}

func createMovieHandleFunc(logger log.Logger, client backend.PageFetcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		level.Debug(logger).Log("msg", "new request")
		amazonID := vars["amazon_id"]

		payload, backErr := client.GetPage("https://www.amazon.de/gp/product/" + amazonID)
		if backErr != nil && backErr.StatusCode == 404 {
			writeErrorReply(logger, w, 404, "the requested movie is not found")
			return
		}
		if backErr != nil {
			level.Error(logger).Log("msg", "backend connection", "error", backErr.Error())
			writeErrorReply(logger, w, backErr.StatusCode, "internal server error")
			return
		}
		defer payload.Close()

		parser := parse.AmazonPrimeParser{}
		parsedPayload, err := parser.Parse(payload)
		if err != nil {
			level.Error(logger).Log("msg", "parsing", "error", err)
			writeErrorReply(logger, w, 500, "internal server error")
			return
		}
		jsonString, err := json.Marshal(parsedPayload)
		if err != nil {
			level.Error(logger).Log("msg", "decoding reply", "error", err)
			writeErrorReply(logger, w, 500, "internal server error")
			return
		}
		writeReply(logger, w, 200, string(jsonString))
	}
}

func writeReply(logger log.Logger, w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write([]byte(message + "\n")); err != nil {
		level.Error(logger).Log("msg", "writing connection", "error", err)
	}
}

func writeErrorReply(logger log.Logger, w http.ResponseWriter, statusCode int, message string) {
	jsonString, err := json.Marshal(ErrorReply{Error: message})
	if err != nil {
		level.Error(logger).Log("msg", "decoding reply", "error", err)
		return
	}
	writeReply(logger, w, statusCode, string(jsonString))
}
