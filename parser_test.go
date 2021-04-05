package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestHtmlTraverse(t *testing.T) {

	html := `
<html>
	<head>
		<title>Page Title</title>
	</head>
	<body>
		<div></div>
		<div>
			<ul>
				<li>first</li>
				<li>second</li>
			</ul>
		</div>
		<p>this is some text</p>
	</body>
</html>`
	input := strings.NewReader(html)
	parsingMap := map[string][]HtmlNode{
		"text":  {{Branch: 2, Tag: "body"}, {Branch: 3, Tag: "p"}},
		"items": {{Branch: 2, Tag: "div"}, {Branch: 1, Tag: "ul"}, {Branch: 0, Tag: "li"}},
	}
	parsedTokens, err := HtmlTraverse(input, parsingMap)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(parsedTokens) != 2 {
		t.Errorf("unexpected parsed tokens length")
	}
	if fmt.Sprintf("%v", parsedTokens["items"]) != "[<li> <li>]" {
		t.Errorf("unexpected parsed tokens[items]")
	}
	if fmt.Sprintf("%v", parsedTokens["text"]) != "[<p>]" {
		t.Errorf("unexpected parsed tokens[text]")
	}
}

func TestHtmlTraverseBranchSelection(t *testing.T) {
	html := `
<html>
	<head>
		<title>Page Title</title>
	</head>
	<body>
		<div></div>
		<div>
			<ul>
				<li>first</li>
				<li>second</li>
			</ul>
		</div>
		<p>this is some text</p>
	</body>
</html>`
	input := strings.NewReader(html)
	parsingMap := map[string][]HtmlNode{
		"onlySecondItem": {{Branch: 2, Tag: "div"}, {Branch: 1, Tag: "ul"}, {Branch: 2, Tag: "li"}},
	}
	parsedTokens, err := HtmlTraverse(input, parsingMap)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(parsedTokens) != 1 {
		t.Errorf("unexpected parsed tokens length")
	}
	if fmt.Sprintf("%v", parsedTokens["onlySecondItem"]) != "[<li>]" {
		t.Errorf("unexpected parsed tokens[onlySecondItem]")
	}
}
