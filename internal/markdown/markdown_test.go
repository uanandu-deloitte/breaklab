package markdown

import "testing"

func TestRenderHeadings(t *testing.T) {
	cases := []struct {
		src  string
		want string
	}{
		{"# Title", "<h1>Title</h1>\n"},
		{"## Section", "<h2>Section</h2>\n"},
		{"### Subsection", "<h3>Subsection</h3>\n"},
	}
	for _, c := range cases {
		if got := string(Render(c.src)); got != c.want {
			t.Errorf("Render(%q) = %q, want %q", c.src, got, c.want)
		}
	}
}

func TestRenderParagraph(t *testing.T) {
	got := string(Render("First line\nsecond line.\n\nSecond paragraph."))
	want := "<p>First line second line.</p>\n<p>Second paragraph.</p>\n"
	if got != want {
		t.Errorf("Render = %q, want %q", got, want)
	}
}

func TestRenderList(t *testing.T) {
	got := string(Render("- one\n- two\n- three"))
	want := "<ul>\n<li>one</li>\n<li>two</li>\n<li>three</li>\n</ul>\n"
	if got != want {
		t.Errorf("Render = %q, want %q", got, want)
	}
}

func TestRenderFencedCode(t *testing.T) {
	got := string(Render("```bash\necho hi\n```"))
	want := "<pre><code>echo hi</code></pre>\n"
	if got != want {
		t.Errorf("Render = %q, want %q", got, want)
	}
}

func TestRenderInlineEmphasisAndCode(t *testing.T) {
	got := string(Render("Use **bold**, *italic*, and `inline code`."))
	want := "<p>Use <strong>bold</strong>, <em>italic</em>, and <code>inline code</code>.</p>\n"
	if got != want {
		t.Errorf("Render = %q, want %q", got, want)
	}
}

func TestRenderEscapesRawHTML(t *testing.T) {
	got := string(Render("<script>alert(1)</script> & \"quotes\""))
	if got == "" {
		t.Fatal("empty output")
	}
	for _, bad := range []string{"<script>", "</script>"} {
		if contains(got, bad) {
			t.Errorf("Render output contains unescaped %q: %s", bad, got)
		}
	}
	for _, want := range []string{"&lt;script&gt;", "&amp;"} {
		if !contains(got, want) {
			t.Errorf("Render output missing escaped %q: %s", want, got)
		}
	}
}

func TestRenderMixedBlocks(t *testing.T) {
	src := "# Title\n\nAn intro paragraph.\n\n- step one\n- step two\n\n```\nverbatim\n```\n\n## Next\n\nClosing line."
	got := string(Render(src))
	for _, want := range []string{
		"<h1>Title</h1>", "<p>An intro paragraph.</p>", "<ul>\n<li>step one</li>",
		"<pre><code>verbatim</code></pre>", "<h2>Next</h2>", "<p>Closing line.</p>",
	} {
		if !contains(got, want) {
			t.Errorf("Render output missing %q:\n%s", want, got)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})()
}
