package parse

import (
	"golang.org/x/net/html"
	"io"
	"strconv"
)

type AmazonPrimePage struct {
	Title       string   `json:"title"`
	ReleaseYear int      `json:"release_year"`
	Actors      []string `json:"actors"`
	Poster      string   `json:"poster"`
	SimilarIDs  []string `json:"similar_ids"`
}

type AmazonPrimeParser struct {
}

var _ Parser = AmazonPrimeParser{}

func (a AmazonPrimeParser) Parse(input io.Reader) (interface{}, error) {
	parser := NewHtmlParser()
	parser.Capture("title", []HtmlNode{{Branch: 2, Tag: "div"}, {Branch: 0, Tag: "div"}, {Branch: 1, Tag: "div"}, {Branch: 1, Tag: "h1"}, {Branch: 1, Tag: "text"}})
	parser.Capture("release_year", []HtmlNode{{Tag: "span", Attributes: map[string]string{"data-automation-id": "release-year-badge"}}, {Tag: "text"}})
	parser.Capture("actor", []HtmlNode{{Branch: 1, Tag: "div"}, {Branch: 4, Tag: "div"}, {Branch: 1, Tag: "div"}, {Branch: 1, Tag: "div"}, {Branch: 2, Tag: "dl"}, {Branch: 2, Tag: "dd"}, {Branch: 0, Tag: "a"}, {Branch: 1, Tag: "text"}})
	parser.Capture("similar_ids", []HtmlNode{{Branch: 1, Tag: "ul"}, {Branch: 0, Tag: "li"}, {Branch: 1, Tag: "div"}, {Branch: 1, Tag: "div"}, {Branch: 1, Tag: "a"}})
	parser.Capture("poster", []HtmlNode{{Branch: 2, Tag: "div"}, {Branch: 3, Tag: "img"}})

	err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	page := newAmazonPrimePage(parser.CapturedNodes)
	return page, nil
}

func newAmazonPrimePage(parsedNodes map[string][]HtmlNode) AmazonPrimePage {
	page := AmazonPrimePage{
		Actors:     make([]string, 0),
		SimilarIDs: make([]string, 0),
	}
	for k, v := range parsedNodes {
		for _, node := range v {
			switch k {
			case "title":
				page.Title = node.Token.Data
			case "release_year":
				year, _ := strconv.ParseInt(node.Token.Data, 10, 32)
				page.ReleaseYear = int(year)
			case "actor":
				page.Actors = append(page.Actors, node.Token.Data)
			case "similar_ids":
				similarID := parseAmazonID(getTokenAttribute(node.Token.Attr, "href"))
				page.SimilarIDs = append(page.SimilarIDs, similarID)
			case "poster":
				page.Poster = getTokenAttribute(node.Token.Attr, "src")
			}
		}
	}
	return page
}

func parseAmazonID(url string) string {
	chars := []byte(url)
	return string(chars[17:27])
}

func getTokenAttribute(attrs []html.Attribute, key string) string {
	for _, attr := range attrs {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
