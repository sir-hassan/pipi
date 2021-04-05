package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"golang.org/x/net/html"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const apiPort = 8080

type Reply struct {
	Title       string   `json:"title"`
	ReleaseYear int      `json:"release_year"`
	Actors      []string `json:"actors"`
	Poster      string   `json:"poster"`
	SimilarIDs  []string `json:"similar_ids"`
}

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
		level.Debug(logger).Log("msg", "new request")
		amazonID := vars["amazon_id"]

		req, _ := http.NewRequest("GET", "https://www.amazon.de/gp/product/"+amazonID, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:86.0) Gecko/20100101 Firefox/86.0")
		httpClient := &http.Client{}
		res, err := httpClient.Do(req)
		if err != nil {
			level.Error(logger).Log("msg", "downstream connection", "error", err)
			writeReply(logger, w, 500, "internal server error")
			return
		}
		defer res.Body.Close()

		parsingMap := map[string][]HtmlNode{
			"title":        {{Branch: 1, Tag: "div"}, {Branch: 2, Tag: "div"}, {Branch: 2, Tag: "div"}, {Branch: 1, Tag: "div"}, {Branch: 1, Tag: "h1"}, {Branch: 1, Tag: "text"}},
			"release_year": {{Branch: 1, Tag: "div"}, {Branch: 2, Tag: "div"}, {Branch: 4, Tag: "span"}, {Branch: 1, Tag: "span"}, {Branch: 1, Tag: "text"}},
			"actor":        {{Branch: 1, Tag: "div"}, {Branch: 4, Tag: "div"}, {Branch: 1, Tag: "div"}, {Branch: 1, Tag: "div"}, {Branch: 2, Tag: "dl"}, {Branch: 2, Tag: "dd"}, {Branch: 0, Tag: "a"}, {Branch: 1, Tag: "text"}},
			"similar_ids":  {{Branch: 1, Tag: "ul"}, {Branch: 0, Tag: "li"}, {Branch: 1, Tag: "div"}, {Branch: 1, Tag: "div"}, {Branch: 1, Tag: "a"}},
			"poster":       {{Branch: 1, Tag: "div"}, {Branch: 2, Tag: "div"}, {Branch: 3, Tag: "img"}},
		}
		parsedTokens, err := HtmlTraverse(res.Body, parsingMap)
		if err != nil {
			level.Error(logger).Log("msg", "parsing", "error", err)
			writeReply(logger, w, 500, "internal server error")
			return
		}
		reply := newReply(parsedTokens)

		jsonString, err := json.Marshal(reply)
		if err != nil {
			level.Error(logger).Log("msg", "decoding reply", "error", err)
			writeReply(logger, w, 500, "internal server error")
			return
		}
		writeReply(logger, w, 200, string(jsonString))
	}
}

func writeReply(logger log.Logger, w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	if _, err := w.Write([]byte(message)); err != nil {
		level.Error(logger).Log("msg", "writing connection", "error", err)
	}
}

func newReply(parsedTokens map[string][]html.Token) Reply {
	reply := Reply{
		Actors:     make([]string, 0),
		SimilarIDs: make([]string, 0),
	}
	for k, v := range parsedTokens {
		for _, token := range v {
			switch k {
			case "title":
				reply.Title = token.Data
			case "release_year":
				year, _ := strconv.ParseInt(token.Data, 10, 32)
				reply.ReleaseYear = int(year)
			case "actor":
				reply.Actors = append(reply.Actors, token.Data)
			case "similar_ids":
				similarID := parseAmazonID(getAttr(token.Attr, "href"))
				reply.SimilarIDs = append(reply.SimilarIDs, similarID)
			case "poster":
				reply.Poster = getAttr(token.Attr, "src")
			}
		}
	}
	return reply
}

func getAttr(attrs []html.Attribute, key string) string {
	for _, attr := range attrs {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func parseAmazonID(url string) string {
	chars := []byte(url)
	return string(chars[17:27])
}
