package main

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

type HtmlNode struct {
	Branch int
	Tag    string
	Token  html.Token
}

func HtmlTraverse(input io.Reader, searchingMap map[string][]HtmlNode) (map[string][]html.Token, error) {
	lastBranch := 0
	currentPath := make([]HtmlNode, 0)

	parsedTokens := make(map[string][]html.Token)
	for k, _ := range searchingMap {
		parsedTokens[k] = make([]html.Token, 0)
	}

	z := html.NewTokenizer(input)
	for {
		tokenType := z.Next()
		token := z.Token()

		switch {
		case tokenType == html.ErrorToken && z.Err() == io.EOF:
			return parsedTokens, nil
		case tokenType == html.ErrorToken:
			return nil, z.Err()
		case tokenType == html.StartTagToken || tokenType == html.TextToken || tokenType == html.SelfClosingTagToken:
			if tokenType == html.TextToken && strings.TrimSpace(token.Data) == "" {
				continue
			}
			currentPath = append(currentPath, HtmlNode{Branch: lastBranch + 1, Token: token})
			lastBranch = 0

			if key := matchSearchingMap(currentPath, searchingMap); key != "" {
				parsedTokens[key] = append(parsedTokens[key], token)
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

func matchPath(path []HtmlNode, searchingPath []HtmlNode) bool {
	for i := 1; i <= len(searchingPath) && i <= len(path); i++ {
		token := path[len(path)-i]
		tag := searchingPath[len(searchingPath)-i]
		if tag.Branch != token.Branch && tag.Branch != 0 {
			return false
		}
		if token.Token.Data != tag.Tag && (tag.Tag != "text" || token.Token.Type != html.TextToken) {
			return false
		}
	}
	return true
}

func matchSearchingMap(path []HtmlNode, searchingMap map[string][]HtmlNode) string {
	for k, v := range searchingMap {
		if matchPath(path, v) {
			return k
		}
	}
	return ""
}
