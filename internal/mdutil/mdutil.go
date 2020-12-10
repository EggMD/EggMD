package mdutil

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

const UNTITLED = "Untitled"

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
