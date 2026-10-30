// Package markdown renders the restricted subset of Markdown a lab's
// lesson.md uses into HTML: headings, paragraphs, unordered lists,
// fenced code blocks, and inline code/bold/italic. It is deliberately
// not a general-purpose Markdown implementation — lesson.md files are
// short, repo-authored concept primers, not arbitrary user content —
// and exists so this repo doesn't need a new external dependency for
// something this small (go.mod carries exactly one today:
// gopkg.in/yaml.v3).
package markdown

import (
	"fmt"
	"html"
	"html/template"
	"regexp"
	"strings"
)

var (
	codeSpanRe = regexp.MustCompile("`([^`]+)`")
	boldRe     = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	italicRe   = regexp.MustCompile(`\*([^*]+)\*`)
)

// Render converts src into HTML. Every raw text run is
// html.EscapeString-ed before any markup is added around it, so content
// can never inject markup of its own — defense in depth even though
// every lesson.md this renders is repo-authored, trusted content.
func Render(src string) template.HTML {
	var out strings.Builder
	lines := strings.Split(src, "\n")

	var para []string
	var list []string

	flushPara := func() {
		if len(para) == 0 {
			return
		}
		out.WriteString("<p>")
		out.WriteString(renderInline(strings.Join(para, " ")))
		out.WriteString("</p>\n")
		para = nil
	}
	flushList := func() {
		if len(list) == 0 {
			return
		}
		out.WriteString("<ul>\n")
		for _, item := range list {
			out.WriteString("<li>")
			out.WriteString(renderInline(item))
			out.WriteString("</li>\n")
		}
		out.WriteString("</ul>\n")
		list = nil
	}

	for i := 0; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])

		if strings.HasPrefix(trimmed, "```") {
			flushPara()
			flushList()
			var code []string
			for i++; i < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[i]), "```"); i++ {
				code = append(code, lines[i])
			}
			out.WriteString("<pre><code>")
			out.WriteString(html.EscapeString(strings.Join(code, "\n")))
			out.WriteString("</code></pre>\n")
			continue
		}

		if text, level, ok := heading(trimmed); ok {
			flushPara()
			flushList()
			fmt.Fprintf(&out, "<h%d>%s</h%d>\n", level, renderInline(text), level)
			continue
		}

		if item, ok := listItem(trimmed); ok {
			flushPara()
			list = append(list, item)
			continue
		}

		if trimmed == "" {
			flushPara()
			flushList()
			continue
		}

		flushList()
		para = append(para, trimmed)
	}
	flushPara()
	flushList()

	return template.HTML(out.String())
}

// heading reports whether line is a "#"/"##"/"###" heading, returning
// its text and level. Checked longest-prefix-first so "### " isn't
// mistaken for "# " followed by literal "## " text.
func heading(line string) (text string, level int, ok bool) {
	for level = 3; level >= 1; level-- {
		prefix := strings.Repeat("#", level) + " "
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(line[len(prefix):]), level, true
		}
	}
	return "", 0, false
}

// listItem reports whether line is a "- " unordered list item.
func listItem(line string) (string, bool) {
	if strings.HasPrefix(line, "- ") {
		return strings.TrimSpace(line[2:]), true
	}
	return "", false
}

// renderInline escapes text and then layers inline code, bold, and
// italic markup on top, in that order — code spans are resolved first
// so a literal "*" inside “ `backticks` “ is never mistaken for
// emphasis, and bold resolves before italic so "**x**" isn't read as
// italic "*" followed by literal "*x*".
func renderInline(text string) string {
	escaped := html.EscapeString(text)
	escaped = codeSpanRe.ReplaceAllString(escaped, "<code>$1</code>")
	escaped = boldRe.ReplaceAllString(escaped, "<strong>$1</strong>")
	escaped = italicRe.ReplaceAllString(escaped, "<em>$1</em>")
	return escaped
}
