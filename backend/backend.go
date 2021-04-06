package backend

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type Error struct {
	StatusCode int
	Err        error
}

func (e Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("status: %d", e.StatusCode)
	}
	return fmt.Sprintf("status: %d, err: %s", e.StatusCode, e.Err)
}

type Client interface {
	GetPage(url string) (io.ReadCloser, *Error)
}

type httpClient struct {
	client *http.Client
}

var _ Client = &httpClient{}

func NewHttpClient(client *http.Client) Client {
	return &httpClient{client: client}
}

func (h httpClient) GetPage(url string) (io.ReadCloser, *Error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:86.0) Gecko/20100101 Firefox/86.0")

	res, err := h.client.Do(req)
	if err != nil {
		return nil, &Error{StatusCode: 500, Err: err}
	}
	if res.StatusCode != 200 {
		return nil, &Error{StatusCode: res.StatusCode}
	}
	return res.Body, nil
}

type filesClient struct {
	rootPath string
}

var _ Client = &filesClient{}

func NewFilesClient(rootPath string) Client {
	return &filesClient{rootPath: rootPath}
}

func (c filesClient) GetPage(url string) (io.ReadCloser, *Error) {
	f, err := os.Open(c.rootPath + "/" + string([]byte(url)[len(url)-10:]))
	if err != nil {
		return nil, &Error{StatusCode: 500, Err: err}
	}
	return f, nil
}
