package gold

import (
	"io"

	"github.com/alecthomas/chroma/quick"
)

type CodeBlockElement struct {
	Code     string
	Language string
}

func (e *CodeBlockElement) Render(w io.Writer, ctx RenderContext) error {
	bs := ctx.blockStack

	var indent uint
	var margin uint
	rules := ctx.style[CodeBlock]
	if rules.Indent != nil {
		indent = *rules.Indent
	}
	if rules.Margin != nil {
		margin = *rules.Margin
	}
	theme := rules.Theme

	iw := &IndentWriter{
		Indent: indent + margin,
		IndentFunc: func(wr io.Writer) {
			renderText(w, bs.Current().Style, " ")
		},
		Forward: &AnsiWriter{
			Forward: w,
		},
	}

	if len(theme) > 0 {
		renderText(iw, bs.Current().Style, rules.Prefix)
		err := quick.Highlight(iw, e.Code, e.Language, "terminal16m", theme)
		renderText(iw, bs.Current().Style, rules.Suffix)
		return err
	}

	// fallback rendering
	el := &BaseElement{
		Token: e.Code,
		Style: rules,
	}

	return el.Render(iw, ctx)
}
