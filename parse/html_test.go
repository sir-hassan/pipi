package parse

import (
	"fmt"
	"strings"
	"testing"
)

func TestHtmlParser(t *testing.T) {

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
	parser := NewHTMLParser()
	parser.Capture("text", []HTMLNode{{Branch: 2, Tag: "body"}, {Branch: 3, Tag: "p"}})
	parser.Capture("items", []HTMLNode{{Branch: 2, Tag: "div"}, {Branch: 1, Tag: "ul"}, {Branch: 0, Tag: "li"}})

	err := parser.Parse(input)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	capturedNodes := parser.CapturedNodes
	if len(capturedNodes) != 2 {
		t.Errorf("unexpected nodes length")
	}
	if fmt.Sprintf("%v", printNodesList(capturedNodes["items"])) != "<li> <li>" {
		t.Errorf("unexpected nodes[items] " + printNodesList(capturedNodes["items"]))
	}
	if fmt.Sprintf("%v", printNodesList(capturedNodes["text"])) != "<p>" {
		t.Errorf("unexpected nodes[text]")
	}
}

func TestHtmlParserBranchSelection(t *testing.T) {
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
	parser := NewHTMLParser()
	parser.Capture("onlySecondItem", []HTMLNode{{Branch: 2, Tag: "div"}, {Branch: 1, Tag: "ul"}, {Branch: 2, Tag: "li"}})

	err := parser.Parse(input)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	capturedNodes := parser.CapturedNodes
	if len(capturedNodes) != 1 {
		t.Errorf("unexpected nodes length")
	}
	if fmt.Sprintf("%v", printNodesList(capturedNodes["onlySecondItem"])) != "<li>" {
		t.Errorf("unexpected nodes[onlySecondItem]")
	}
}

func TestHtmlParserMatchAttr(t *testing.T) {
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
				<li data-foo="bar">second</li>
			</ul>
		</div>
		<p>this is some text</p>
	</body>
</html>`
	input := strings.NewReader(html)
	parser := NewHTMLParser()
	parser.Capture("item", []HTMLNode{{Tag: "ul"}, {Tag: "li", Attributes: map[string]string{"data-foo": "bar"}}})

	err := parser.Parse(input)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	capturedNodes := parser.CapturedNodes
	if len(capturedNodes) != 1 {
		t.Errorf("unexpected nodes length")
	}
	if fmt.Sprintf("%v", printNodesList(capturedNodes["item"])) != "<li data-foo=\"bar\">" {
		t.Errorf("unexpected nodes[item]: %s", printNodesList(capturedNodes["item"]))
	}
}

func printNodesList(nodesList []HTMLNode) string {
	result := ""
	for _, v := range nodesList {
		result += fmt.Sprintf("%v ", v.Token)
	}
	return strings.TrimSpace(result)
}
