package main

import (
	"golang.org/x/net/html"
	"io"
)

type HtmlNode struct {
	Branch int
	Tag    string
}

type Tag struct {
	Branch int
	Token  html.Token
}

func HtmlTraverse(input io.Reader, parsingMap map[string][]HtmlNode) (map[string][]html.Token, error) {
	lastBranch := 0

	sequence := make([]Tag, 0)

	parsedTokens := make(map[string][]html.Token)
	for k, _ := range parsingMap {
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
			sequence = append(sequence, Tag{Branch: lastBranch + 1, Token: token})
			lastBranch = 0

			if key := matchPath(sequence, parsingMap); key != "" {
				parsedTokens[key] = append(parsedTokens[key], token)
			}

			//if token.Data == "img" && hasAttr(token.Attr, "src", "https://images-na.ssl-images-amazon.com/images/S/sgp-catalog-images/region_DE/universum-00664000-Full-Image_GalleryBackground-de-DE-1617099345129._SX1080_.jpg") {
			//	//fmt.Printf("We found a link!: %v\n", sequence)
			//}

			if tokenType == html.TextToken || tokenType == html.SelfClosingTagToken {
				lastTag := sequence[len(sequence)-1]
				lastBranch = lastTag.Branch
				sequence = sequence[:len(sequence)-1]
			}

		case tokenType == html.EndTagToken:
			lastTag := sequence[len(sequence)-1]
			lastBranch = lastTag.Branch
			sequence = sequence[:len(sequence)-1]
		}
	}
}

func match(sequence []Tag, path []HtmlNode) bool {
	for i := 1; i <= len(path) && i <= len(sequence); i++ {
		token := sequence[len(sequence)-i]
		tag := path[len(path)-i]
		if tag.Branch != token.Branch && tag.Branch != 0 {
			return false
		}
		if token.Token.Data != tag.Tag && (tag.Tag != "text" || token.Token.Type != html.TextToken) {
			return false
		}
	}
	return true
}

func matchPath(sequence []Tag, parsingMap map[string][]HtmlNode) string {
	for k, v := range parsingMap {
		if match(sequence, v) {
			return k
		}
	}
	return ""
}
