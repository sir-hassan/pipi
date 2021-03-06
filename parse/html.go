// `html.Tokenizer` walks through html input text and convert it to sequence on
// tokens. These tokens reflect their respective html elements and encode their
// attributes, type and text.
// We can leverage `html.Tokenizer` to build a simple html searching and
// parsing logic by just looping though these tokens, keeping track of their
// paths to the root element(<body/>) and comparing them with a given searching
// paths.

package parse

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

// HTMLNode can represent an html searching node (used in building searching
// paths) or an html actual node (used in representing parsed html nodes).
type HTMLNode struct {
	// Tag, Branch and Attributes fields are used for searching matching nodes.
	Tag        string
	Branch     int
	Attributes map[string]string

	// Token field represent an actual html token that maps to matching fields.
	Token html.Token
}

// HTMLParser leverages html.Tokenizer and implements simple html dom parsing
// logic. You give searching details of html nodes to capture (the dom path, tag
// type and attributes) then HTMLParser will loop though the document and try
// to capture the matching nodes.
type HTMLParser struct {
	SearchingMap  map[string][]HTMLNode
	CapturedNodes map[string][]HTMLNode
}

// NewHTMLParser creates a new HTMLParser.
func NewHTMLParser() *HTMLParser {
	return &HTMLParser{
		SearchingMap:  make(map[string][]HTMLNode),
		CapturedNodes: make(map[string][]HTMLNode),
	}
}

// Capture records a capturing(searching) rule.
func (t *HTMLParser) Capture(fieldName string, path []HTMLNode) {
	t.SearchingMap[fieldName] = path
	t.CapturedNodes[fieldName] = make([]HTMLNode, 0)
}

// Parse runs the parsing loop against the SearchingMap field and recorders
// the results in CapturedNodes field.
func (t *HTMLParser) Parse(input io.Reader) error {
	lastBranch := 0
	currentPath := make([]HTMLNode, 0)

	z := html.NewTokenizer(input)
	for {
		tokenType := z.Next()
		token := z.Token()

		switch {
		case tokenType == html.ErrorToken && z.Err() == io.EOF:
			return nil
		case tokenType == html.ErrorToken:
			return z.Err()
		case tokenType == html.StartTagToken || tokenType == html.TextToken || tokenType == html.SelfClosingTagToken:
			if tokenType == html.TextToken && strings.TrimSpace(token.Data) == "" {
				continue
			}
			currentPath = append(currentPath, HTMLNode{Branch: lastBranch + 1, Token: token})
			lastBranch = 0

			if key := matchSearchingMap(currentPath, t.SearchingMap); key != "" {
				t.CapturedNodes[key] = append(t.CapturedNodes[key], HTMLNode{Token: token})
			}
			if tokenType == html.TextToken || tokenType == html.SelfClosingTagToken {
				lastTag := currentPath[len(currentPath)-1]
				lastBranch = lastTag.Branch
				currentPath = currentPath[:len(currentPath)-1]
			}
		case tokenType == html.EndTagToken:
			lastTag := currentPath[len(currentPath)-1]
			lastBranch = lastTag.Branch
			currentPath = currentPath[:len(currentPath)-1]
		}
	}
}

func matchPath(path []HTMLNode, searchingPath []HTMLNode) bool {
	for i := 1; i <= len(searchingPath) && i <= len(path); i++ {
		token := path[len(path)-i]
		tag := searchingPath[len(searchingPath)-i]
		if tag.Branch != token.Branch && tag.Branch != 0 {
			return false
		}
		if token.Token.Data != tag.Tag && (tag.Tag != "text" || token.Token.Type != html.TextToken) {
			return false
		}
		if len(tag.Attributes) > 0 {
			attributes := make(map[string]string)
			for _, attr := range token.Token.Attr {
				attributes[attr.Key] = attr.Val
			}
			for key, value := range tag.Attributes {
				if v, ok := attributes[key]; !ok || v != value {
					return false
				}
			}
		}
	}
	return true
}

func matchSearchingMap(path []HTMLNode, searchingMap map[string][]HTMLNode) string {
	for k, v := range searchingMap {
		if matchPath(path, v) {
			return k
		}
	}
	return ""
}
