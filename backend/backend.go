package backend

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Client is the interface that fetches a web page for a given url.
type Client interface {
	GetPage(url string) (io.ReadCloser, *Error)
}

// Error represents a Client interface error while retrieving a page.
// It wrap status code and error.
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

type httpClient struct {
	client *http.Client
}

var _ Client = &httpClient{}

// NewHTTPClient creates a new httpClient
func NewHTTPClient(client *http.Client) Client {
	return &httpClient{client: client}
}

func (h httpClient) GetPage(url string) (io.ReadCloser, *Error) {
	req, _ := http.NewRequest("GET", url, nil)

	// surprisingly this header is enough to let amazon.de think
	// that you are not a robot.
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

// NewFilesClient creates new filesClient for testing purposes.
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
