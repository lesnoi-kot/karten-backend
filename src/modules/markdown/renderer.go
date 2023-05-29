package markdown

import (
	"io"

	"github.com/russross/blackfriday/v2"
)

type CustomHTMLRenderer struct {
	*blackfriday.HTMLRenderer
}

func (r *CustomHTMLRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	return r.HTMLRenderer.RenderNode(w, node, entering)
}

func NewCustomHTMLRenderer() *CustomHTMLRenderer {
	return &CustomHTMLRenderer{
		HTMLRenderer: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			Flags: blackfriday.CommonHTMLFlags,
		}),
	}
}
