package report_renderers

import (
	"strings"
	"unsafe"
	"wb_logistic_assistant/internal/reports"
)

const telegramMessageLimit = 4000

type TelegramBotRenderMode int

const (
	TelegramBotRenderMarkdown TelegramBotRenderMode = iota
	TelegramBotRenderHTML
	TelegramBotRenderPlain
)

type TelegramBotRenderer struct {
	Mode          TelegramBotRenderMode
	IsTitle       bool
	builder       *strings.Builder
	escapeBuilder *strings.Builder
	messages      []string

	prefix                string
	blockOpen, blockClose string
	styleOpen, styleClose string
	linkOpen              bool
	linkSize              int
}

func (r *TelegramBotRenderer) Render(data *reports.ReportData) ([]string, error) {
	r.builder = &strings.Builder{}
	r.escapeBuilder = &strings.Builder{}
	r.messages = []string{}

	if data.Title != "" {
		r.renderTitle(data.Title)
	}
	if data.Header != nil {
		r.render(data.Header)
	}
	if data.Body != nil {
		r.render(data.Body)
	}

	if r.builder.Len() > 0 {
		r.messages = append(r.messages, r.builder.String())
	}

	return r.messages, nil
}

func (r *TelegramBotRenderer) renderTitle(title string) {
	switch r.Mode {
	case TelegramBotRenderMarkdown:
		r.builder.WriteString("*" + r.escapeMarkdown(title) + "*\n\n")
	case TelegramBotRenderHTML:
		r.builder.WriteString("<b>" + title + "</b>\n\n")
	default:
		r.builder.WriteString(title + "\n\n")
	}
}

func (r *TelegramBotRenderer) render(item *reports.Item) {
	if item == nil {
		return
	}
	switch r.Mode {
	case TelegramBotRenderMarkdown:
		r.renderMarkdown(item)
	case TelegramBotRenderHTML:
		r.renderHTML(item)
	default:
		r.renderPlain(item)
	}
}

func (r *TelegramBotRenderer) renderMarkdown(item *reports.Item) {
	for i, it := range item.Children {
		if it.Block {
			if i != 0 {
				r.writeText("\n", true)
				if r.prefix != "" {
					r.writeText(r.prefix, true)
				}
			}
		} else {
			r.writeText(" ", true)
		}

		if it.HiddenQuote {
			r.blockOpen = "**>"
			r.blockClose = " ||"
			r.prefix = ">"
			r.writeText("**>", true)
		} else if it.Quote {
			r.blockOpen = ""
			r.blockClose = ""
			r.prefix = "> "
			r.writeText("> ", true)
		} else if it.Code {
			r.blockOpen = "```"
			r.blockClose = "```"
			r.writeText("```", true)
		}

		if it.Bold {
			r.styleOpen = "*"
			r.styleClose = "*"
			r.writeText("*", true)
		}
		if it.Link != "" {
			r.linkOpen = true
			r.linkSize = len(it.Link) + len(it.Text) + 6
			r.writeText("[", true)
		}

		text := r.escapeMarkdown(it.Text)
		r.writeText(text, false)

		if it.Link != "" {
			r.writeText("](", true)
			r.writeText(it.Link, false)
			r.writeText(")", true)
			r.linkOpen = false
			r.linkSize = 0
		}
		if it.Bold {
			r.styleOpen = ""
			r.styleClose = ""
			r.writeText("*", true)
		}

		if len(it.Children) > 0 {
			r.renderMarkdown(it)
		}

		if it.HiddenQuote {
			r.blockOpen = ""
			r.blockClose = ""
			r.prefix = ""
			r.writeText(" ||", true)
		} else if it.Code {
			r.blockOpen = ""
			r.blockClose = ""
			r.prefix = ""
			r.writeText("\n```", true)
		}
	}
}

func (r *TelegramBotRenderer) renderHTML(item *reports.Item) {
	for i, it := range item.Children {
		if it.Block {
			if i != 0 {
				r.writeText("\n", true)
			}
		} else {
			r.writeText(" ", true)
		}

		if it.HiddenQuote {
			r.blockOpen = "<blockquote expandable>"
			r.blockClose = "</blockquote>"
			r.writeText("<blockquote expandable>", true)
		} else if it.Quote {
			r.blockOpen = "<blockquote>"
			r.blockClose = "</blockquote>"
			r.writeText("<blockquote>", true)
		} else if it.Code {
			r.blockOpen = "<pre><code>"
			r.blockClose = "</code></pre>"
			r.writeText("<pre><code>", true)
		}

		if it.Bold {
			r.styleOpen = "<b>"
			r.styleClose = "</b>"
			r.writeText("<b>", true)
		}
		if it.Link != "" {
			r.linkOpen = true
			r.linkSize = len(it.Link) + len(it.Text) + 15
			r.writeText(`<a href="`, true)
			r.writeText(it.Link, false)
			r.writeText(`">`, true)
		}

		r.writeText(it.Text, false)

		if it.Link != "" {
			r.linkOpen = false
			r.linkSize = 0
			r.writeText("</a>", true)
		}
		if it.Bold {
			r.writeText("</b>", true)
			r.styleOpen = ""
			r.styleClose = ""
		}

		if len(it.Children) > 0 {
			r.renderHTML(it)
		}

		if it.HiddenQuote {
			r.writeText("</blockquote>", true)
			r.blockOpen = ""
			r.blockClose = ""
		} else if it.Quote {
			r.writeText("</blockquote>", true)
			r.blockOpen = ""
			r.blockClose = ""
		} else if it.Code {
			r.writeText("</code></pre>", true)
			r.blockOpen = ""
			r.blockClose = ""
		}
	}
}

func (r *TelegramBotRenderer) renderPlain(item *reports.Item) {
	for i, it := range item.Children {
		if it.Block && i != 0 {
			r.writeText("\n", true)
		} else {
			r.writeText(" ", true)
		}
		r.writeText(it.Text, false)
		if len(it.Children) > 0 {
			r.renderPlain(it)
		}
	}
}

func (r *TelegramBotRenderer) writeText(text string, tag bool) {
	b := r.builder
	curLen := b.Len()
	textLen := len(text)

	const overhead = 15
	limit := telegramMessageLimit

	if curLen+textLen+overhead < limit {
		b.WriteString(text)
		return
	}

	if r.linkOpen && r.linkSize < limit && curLen+r.linkSize >= limit {
		r.fastFlush()
		r.fastWriteOpenTags()
		b.WriteString(text)
		return
	}

	if tag {
		r.fastFlush()
		r.fastWriteOpenTags()
		b.WriteString(text)
		return
	}

	writeLen := limit - curLen - overhead
	if writeLen > 0 {
		b.WriteString(text[:writeLen])
	}

	r.fastFlush()
	r.fastWriteOpenTags()

	if writeLen < textLen {
		b.WriteString(text[writeLen:])
	}
}

func (r *TelegramBotRenderer) fastFlush() {
	b := r.builder

	if r.styleClose != "" {
		b.WriteString(r.styleClose)
	}
	if r.blockClose != "" {
		b.WriteString(r.blockClose)
	}

	r.messages = append(r.messages, b.String())
	b.Reset()
}

func (r *TelegramBotRenderer) fastWriteOpenTags() {
	b := r.builder
	if r.blockOpen != "" {
		b.WriteString(r.blockOpen)
	}
	if r.styleOpen != "" {
		b.WriteString(r.styleOpen)
	}
}

func (r *TelegramBotRenderer) escapeMarkdown(s string) string {
	n := len(s)
	buf := make([]byte, 0, n+n/8)

	for i := 0; i < n; i++ {
		ch := s[i]
		switch ch {
		case '!', '*', '_', '-', '~', '`', '[', ']', '(', ')', '>', '|', '.':
			buf = append(buf, '\\', ch)
		default:
			buf = append(buf, ch)
		}
	}

	return *(*string)(unsafe.Pointer(&buf))
}
