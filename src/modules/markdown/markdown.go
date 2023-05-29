package markdown

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

var markdownPolicy = bluemonday.UGCPolicy()

func init() {
	markdownPolicy.AddTargetBlankToFullyQualifiedLinks(true)
	markdownPolicy.RequireParseableURLs(true)
	markdownPolicy.RequireNoReferrerOnLinks(true)
	markdownPolicy.RequireSandboxOnIFrame()

	markdownPolicy.AllowAttrs("class").
		Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).
		OnElements("code")
}

func Render(md string) string {
	unsafeHTML := blackfriday.Run(
		[]byte(md),
		blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.AutoHeadingIDs|blackfriday.HardLineBreak),
		blackfriday.WithRenderer(NewCustomHTMLRenderer()),
	)
	html := markdownPolicy.SanitizeBytes(unsafeHTML)
	return string(html)
}
