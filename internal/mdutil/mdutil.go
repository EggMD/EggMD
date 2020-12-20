package mdutil

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

const UNTITLED = "Untitled"

// ParseTitle returns the document title by parsing the giving markdown.
func ParseTitle(content string) string {
	reader := text.NewReader([]byte(content))
	p := goldmark.DefaultParser()
	p.AddOptions(parser.WithAttribute())
	mdAST := p.Parse(reader)
	if mdAST == nil || mdAST.FirstChild() == nil {
		return UNTITLED
	}

	if mdAST.FirstChild().Kind().String() == "Heading" {
		heading, ok := mdAST.FirstChild().(*ast.Heading)
		if !ok {
			return UNTITLED
		}
		txt, ok := heading.FirstChild().(*ast.Text)
		if !ok {
			return UNTITLED
		}

		return content[txt.Segment.Start:txt.Segment.Stop]
	}

	return UNTITLED
}

// RenderMarkdown renders the input markdown string.
func RenderMarkdown(rawMarkdown string) (string, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(rawMarkdown), &buf); err != nil {
		return "", err
	}

	return bluemonday.UGCPolicy().Sanitize(buf.String()), nil
}
