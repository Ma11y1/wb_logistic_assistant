package report_renderers

import (
	"strings"
	"unsafe"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/reports"
)

const (
	telegramMessageLimit         = 4000
	telegramMessageLimitOverhead = 15
)

type TelegramBotRenderMode int

const (
	TelegramBotRenderMarkdown TelegramBotRenderMode = iota
	TelegramBotRenderHTML
	TelegramBotRenderPlain
)

type TelegramBotRenderer struct {
	Mode     TelegramBotRenderMode
	builder  *strings.Builder
	messages []string

	prefix                string
	blockOpen, blockClose string
	styleOpen, styleClose string
	linkOpen              bool
	linkSize              int
}

func (r *TelegramBotRenderer) Render(data *reports.ReportData) ([]string, error) {
	r.builder = &strings.Builder{}
	r.builder.Grow(telegramMessageLimit)
	r.messages = []string{}

	r.prefix = ""
	r.blockOpen, r.blockClose = "", ""
	r.styleOpen, r.styleClose = "", ""
	r.linkOpen, r.linkSize = false, 0

	if data.Header != nil {
		if err := r.render(data.Header); err != nil {
			return nil, err
		} else {
			r.write("\n", true)
		}
	}
	if data.Body != nil {
		if err := r.render(data.Body); err != nil {
			return nil, err
		}
	}

	if r.builder.Len() > 0 {
		r.messages = append(r.messages, r.builder.String())
	}

	return r.messages, nil
}

func (r *TelegramBotRenderer) render(item *reports.Item) error {
	if item == nil {
		return nil
	}
	switch r.Mode {
	case TelegramBotRenderMarkdown:
		return r.renderMarkdown(item)
	case TelegramBotRenderHTML:
		return r.renderHTML(item)
	default:
		r.renderPlain(item)
		return nil
	}
}

// Plain
func (r *TelegramBotRenderer) renderPlain(item *reports.Item) {
	for i, child := range item.Children {
		if child.Block && i != 0 {
			r.write("\n", false)
		} else {
			r.write(" ", false)
		}

		r.write(child.Text, false)
		if len(child.Link) > 0 {
			r.write("["+child.Link+"]", false)
		}

		if len(child.Children) > 0 {
			r.renderPlain(child)
		}
	}
}

// Markdown
func (r *TelegramBotRenderer) renderMarkdown(item *reports.Item) error {
	for i, child := range item.Children {
		if child.Block && i != 0 {
			r.write("\n", true)
			if r.prefix != "" {
				r.write(r.prefix, true)
			}
		} else if i != 0 {
			r.write(" ", true)
		}

		if child.HiddenQuote {
			r.blockOpen, r.blockClose, r.prefix = "**>", " ||", "> "
			r.write(r.blockOpen, true)
		} else if child.Quote {
			r.blockOpen, r.blockClose, r.prefix = "> ", "", "> "
			r.write(r.blockOpen, true)
		} else if child.Code {
			r.blockOpen, r.blockClose = "```", "```"
			r.write(r.blockOpen, true)
		}

		if child.Bold {
			r.styleOpen, r.styleClose = "*", "*"
			r.write(r.styleOpen, true)
		}
		if child.Link != "" {
			// Reflection symbols are placed in advance to accurately determine the link size
			child.Link = r.escapeMarkdown(child.Link)
			child.Text = r.escapeMarkdown(child.Text)
			r.linkOpen = true
			r.linkSize = len(child.Link) + len(child.Text) + 6
			if r.linkSize >= telegramMessageLimit {
				return errors.Newf("TelegramBotRenderer.renderMarkdown", "link is very large, max size %d, link size %d", telegramMessageLimit, r.linkSize)
			}
			r.write("[", true)
		}

		r.write(r.escapeMarkdown(child.Text), false)

		if child.Link != "" {
			// Recalculation of the link size since the link tag was previously opened and the description text was inserted, and the next time the link itself is recorded with the closing of the tag, the remaining size should be correct
			r.linkSize = len(child.Link) + 3
			r.write("]("+child.Link+")", true)
			r.linkOpen = false
			r.linkSize = 0
		}
		if child.Bold {
			r.styleClose = ""
			r.styleOpen = ""
			r.write("*", true)
		}

		if len(child.Children) > 0 {
			if err := r.renderMarkdown(child); err != nil {
				return err
			}
		}

		if child.HiddenQuote {
			r.blockOpen, r.blockClose, r.prefix = "", "", ""
			r.write(" ||", true)
		} else if child.Code {
			r.blockOpen, r.blockClose, r.prefix = "", "", ""
			r.write("\n```", true)
		}
	}
	return nil
}

func (r *TelegramBotRenderer) escapeMarkdown(s string) string {
	if s == "" {
		return ""
	}
	n := len(s)
	buf := make([]byte, 0, n+n/8)
	for i := 0; i < n; i++ {
		ch := s[i]
		switch ch {
		case '!', '*', '_', '-', '~', '`', '[', ']', '(', ')', '>', '|', '.':
			// skip if reflection already exists
			if i != 0 && s[i-1] == '\\' {
				continue
			}
			buf = append(buf, '\\', ch)
		default:
			buf = append(buf, ch)
		}
	}
	return *(*string)(unsafe.Pointer(&buf))
}

// HTML
func (r *TelegramBotRenderer) renderHTML(item *reports.Item) error {
	for i, child := range item.Children {
		if child.Block && i != 0 {
			r.write("\n", true)
		} else if i != 0 {
			r.write(" ", true)
		}

		if child.HiddenQuote {
			r.blockOpen, r.blockClose = "<blockquote expandable>", "</blockquote>"
			r.write(r.blockOpen, true)
		} else if child.Quote {
			r.blockOpen, r.blockClose = "<blockquote>", "</blockquote>"
			r.write(r.blockOpen, true)
		} else if child.Code {
			r.blockOpen, r.blockClose = "<pre><code>", "</code></pre>"
			r.write(r.blockOpen, true)
		}

		if child.Bold {
			r.styleOpen, r.styleClose = "<b>", "</b>"
			r.write(r.styleOpen, true)
		}
		if child.Link != "" {
			r.linkOpen = true
			r.linkSize = len(child.Link) + len(child.Text) + 15
			if r.linkSize >= telegramMessageLimit {
				return errors.Newf("TelegramBotRenderer.renderHTML()", "link is very large, max size %d, link size %d", telegramMessageLimit, r.linkSize)
			}
			r.write(`<a href="`+child.Link+`">`, true)
		}

		r.write(child.Text, false)

		if child.Link != "" {
			r.linkOpen = false
			r.linkSize = 0
			r.write("</a>", true)
		}
		if child.Bold {
			r.write("</b>", true)
			r.styleOpen = ""
			r.styleClose = ""
		}

		if len(child.Children) > 0 {
			err := r.renderHTML(child)
			if err != nil {
				return err
			}
		}

		switch {
		case child.HiddenQuote, child.Quote:
			r.write("</blockquote>", true)
		case child.Code:
			r.write("</code></pre>", true)
		}
	}
	return nil
}

func (r *TelegramBotRenderer) write(text string, tag bool) {
	if text == "" {
		return
	}
	builder := r.builder
	builderLen := builder.Len()
	textLen := len(text)

	// if a link is open and the entire link needs to be placed there, we move the entire link into a new message
	if r.linkOpen && r.linkSize < telegramMessageLimit && builderLen+r.linkSize > telegramMessageLimit {
		r.flush()
		r.writeOpenTags()
		builder.WriteString(text)
		return
	}

	// if it fits completely, write it down immediately.
	if builderLen+textLen+telegramMessageLimitOverhead < telegramMessageLimit {
		builder.WriteString(text)
		return
	}

	if tag {
		r.flush()
		r.writeOpenTags()
		builder.WriteString(text)
		return
	}

	writeLen := telegramMessageLimit - builderLen - telegramMessageLimitOverhead
	if writeLen > 0 {
		safeText := safeCutTextAlongRune(text, writeLen)
		if r.Mode == TelegramBotRenderMarkdown && safeText != "" && safeText[len(safeText)-1] == '\\' {
			text = text[len(safeText)-1:]
			safeText = safeText[:len(safeText)-1]
		} else {
			text = text[len(safeText):]
		}
		if safeText != "" {
			builder.WriteString(safeText)
		}
	} else {
		text = ""
	}

	r.flush()
	r.writeOpenTags()

	if len(text) > 0 {
		builder.WriteString(text)
	}
}

func (r *TelegramBotRenderer) flush() {
	if r.styleOpen != "" {
		r.builder.WriteString(r.styleClose)
	}
	if r.blockOpen != "" {
		r.builder.WriteString(r.blockClose)
	}
	r.messages = append(r.messages, r.builder.String())
	r.builder.Reset()
}

func (r *TelegramBotRenderer) writeOpenTags() {
	if r.blockOpen != "" {
		r.builder.WriteString(r.blockOpen)
	}
	if r.styleOpen != "" {
		r.builder.WriteString(r.styleOpen)
	}
}

func safeCutTextAlongRune(s string, lim int) string {
	if len(s) <= lim {
		return s
	}
	i := lim
	for i > 0 && (s[i]&0xC0) == 0x80 {
		i--
	}
	return s[:i]
}
