package formatter

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FormatNode(sel *goquery.Selection) string {
	tag := goquery.NodeName(sel)

	if tag == "img" {
		src, _ := sel.Attr("src")
		return fmt.Sprintf("<img src='%s' class='mt-4 p-2 max-w-full h-auto'>", src)
	}

	var contentBuilder strings.Builder
	sel.Contents().Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == "#text" {
			contentBuilder.WriteString(s.Text())
		} else {
			contentBuilder.WriteString(FormatNode(s))
		}
	})
	content := contentBuilder.String()

	switch tag {
	case "p":
		return fmt.Sprintf("<p class='text-lg mt-4'>%s</p>", content)
	case "h2":
		return fmt.Sprintf("<h2 class='text-2xl mt-4 p-2'>%s</h2>", content)
	case "pre":
		return fmt.Sprintf("<pre data-prefix=''>%s</pre>", content)
	case "blockquote":
		return fmt.Sprintf("<blockquote class='text-lg mt-4 p-2'>%s</blockquote>", content)
	default:
		if tag == "" || tag == "#text" {
			return content
		}
		return fmt.Sprintf("<%s>%s</%s>", tag, content, tag)
	}
}
