package formatter

import (
	"fmt"
	"html"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/*
This function takes a goquery.Selection and returns a formatted string of the HTML content.

It does so by matching the tag of the selection and then calling processNodeContents to recursively
format the contents of the selection.
*/

func FormatNode(sel *goquery.Selection) string {
	if sel == nil {
		return ""
	}

	tag := goquery.NodeName(sel)

	switch {
	case tag == "#text":
		return html.EscapeString(sel.Text())

	case tag == "img":
		src, _ := sel.Attr("src")
		alt, _ := sel.Attr("alt")
		return fmt.Sprintf("<img src='%s' alt='%s' class='mt-4 p-2 max-w-full h-auto rounded shadow-sm'>", src, alt)

	case tag == "a":
		href, _ := sel.Attr("href")
		content := processNodeContents(sel)
		return fmt.Sprintf("<a href='%s' class='link-primary hover:underline'>%s</a>", href, content)

	case tag == "code":
		content := processNodeContents(sel)
		return fmt.Sprintf("<code class='px-1 py-0.5 rounded text-sm font-mono'>%s</code>", content)

	case tag == "ul":
		content := processNodeContents(sel)
		return fmt.Sprintf("<ul class='list-disc pl-5 mt-4 space-y-2'>%s</ul>", content)

	case tag == "ol":
		content := processNodeContents(sel)
		return fmt.Sprintf("<ol class='list-decimal pl-5 mt-4 space-y-2'>%s</ol>", content)

	case tag == "li":
		content := processNodeContents(sel)
		return fmt.Sprintf("<li>%s</li>", content)

	case tag == "p":
		content := processNodeContents(sel)
		return fmt.Sprintf("<p class='text-lg mt-4'>%s</p>", content)

	case tag == "h1":
		content := processNodeContents(sel)
		return fmt.Sprintf("<h1 class='text-3xl font-bold mt-6 mb-2'>%s</h1>", content)

	case tag == "h2":
		content := processNodeContents(sel)
		return fmt.Sprintf("<h2 class='text-2xl font-semibold mt-5 mb-2'>%s</h2>", content)

	case tag == "h3":
		content := processNodeContents(sel)
		return fmt.Sprintf("<h3 class='text-xl font-medium mt-4 mb-2'>%s</h3>", content)

	case tag == "h4", tag == "h5", tag == "h6":
		content := processNodeContents(sel)
		return fmt.Sprintf("<%s class='font-medium mt-3 mb-2'>%s</%s>", tag, content, tag)

	case tag == "pre":
		content := processNodeContents(sel)
		return fmt.Sprintf("<pre class='p-4 rounded overflow-x-auto font-mono text-sm mt-4' data-prefix=''>%s</pre>", content)

	case tag == "blockquote":
		content := processNodeContents(sel)
		return fmt.Sprintf("<blockquote class='border-l-4 border-gray-300 pl-4 italic mt-4'>%s</blockquote>", content)

	case tag == "table":
		content := processNodeContents(sel)
		return fmt.Sprintf("<div class='overflow-x-auto mt-4'><table class='min-w-full divide-y divide-gray-200'>%s</table></div>", content)

	case tag == "tr":
		content := processNodeContents(sel)
		return fmt.Sprintf("<tr>%s</tr>", content)

	case tag == "th":
		content := processNodeContents(sel)
		return fmt.Sprintf("<th class='px-6 py-3 text-left text-xs font-medium uppercase tracking-wider'>%s</th>", content)

	case tag == "td":
		content := processNodeContents(sel)
		return fmt.Sprintf("<td class='px-6 py-4 whitespace-nowrap'>%s</td>", content)

	case tag == "strong", tag == "b":
		content := processNodeContents(sel)
		return fmt.Sprintf("<strong class='font-bold'>%s</strong>", content)

	case tag == "em", tag == "i":
		content := processNodeContents(sel)
		return fmt.Sprintf("<em class='italic'>%s</em>", content)

	case tag == "span":
		content := processNodeContents(sel)
		return fmt.Sprintf("<span>%s</span>", content)

	case tag == "div":
		content := processNodeContents(sel)
		return fmt.Sprintf("<div>%s</div>", content)

	case tag == "hr":
		return "<hr class='my-4 border-t border-gray-300'>"

	case tag == "br":
		return "<br>"

	case tag == "iframe":
		src, _ := sel.Attr("src")
		width, _ := sel.Attr("width")
		height, _ := sel.Attr("height")
		return fmt.Sprintf("<iframe src='%s' width='%s' height='%s' class='border-0 mt-4' allowfullscreen></iframe>", src, width, height)

	case tag == "":
		return processNodeContents(sel)

	default:
		content := processNodeContents(sel)
		class, hasClass := sel.Attr("class")
		id, hasID := sel.Attr("id")

		attributes := ""
		if hasClass {
			attributes += fmt.Sprintf(" class='%s'", class)
		}
		if hasID {
			attributes += fmt.Sprintf(" id='%s'", id)
		}

		return fmt.Sprintf("<%s%s>%s</%s>", tag, attributes, content, tag)
	}
}

func processNodeContents(sel *goquery.Selection) string {
	var contentBuilder strings.Builder
	sel.Contents().Each(func(i int, s *goquery.Selection) {
		contentBuilder.WriteString(FormatNode(s))
	})
	return contentBuilder.String()
}

func getAttributes(sel *goquery.Selection) string {
	var attrBuilder strings.Builder

	for _, attr := range sel.Nodes[0].Attr {
		if attr.Key == "style" || attr.Key == "onclick" {
			continue
		}

		escapedValue := html.EscapeString(attr.Val)
		attrBuilder.WriteString(fmt.Sprintf(" %s='%s'", attr.Key, escapedValue))
	}

	return attrBuilder.String()
}
