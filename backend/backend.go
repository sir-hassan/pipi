package backend

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// PageFetcher is the interface that fetches a web page for a given url. PageFetcher
// should be safe for concurrent use by multiple goroutines and implementations
// should respect that.
type PageFetcher interface {
	GetPage(url string) (io.ReadCloser, *Error)
}

// Error represents a PageFetcher interface error while retrieving a page.
// It wraps status code and error.
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

type webPageFetcher struct {
	client *http.Client
}

var _ PageFetcher = &webPageFetcher{}

// NewWebPageFetcher creates a new webPageFetcher
func NewWebPageFetcher(client *http.Client) PageFetcher {
	return &webPageFetcher{client: client}
}

func (p webPageFetcher) GetPage(url string) (io.ReadCloser, *Error) {
	req, _ := http.NewRequest("GET", url, nil)

	// surprisingly this header is enough to let amazon.de think
	// that you are not a robot.
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:86.0) Gecko/20100101 Firefox/86.0")

	res, err := p.client.Do(req)
	if err != nil {
		return nil, &Error{StatusCode: 500, Err: err}
	}
	if res.StatusCode != 200 {
		return nil, &Error{StatusCode: res.StatusCode}
	}
	return res.Body, nil
}

type localPageFetcher struct {
	rootPath string
}

var _ PageFetcher = &localPageFetcher{}

// NewLocalPageFetcher creates new localPageFetcher for testing purposes.
func NewLocalPageFetcher(rootPath string) PageFetcher {
	return &localPageFetcher{rootPath: rootPath}
}

func (p localPageFetcher) GetPage(url string) (io.ReadCloser, *Error) {
	f, err := os.Open(p.rootPath + "/" + string([]byte(url)[len(url)-10:]))
	if err != nil {
		return nil, &Error{StatusCode: 500, Err: err}
	}
	return f, nil
}
