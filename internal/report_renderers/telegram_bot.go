package report_renderers

import (
	"strings"
	"wb_logistic_assistant/internal/logger"
	"wb_logistic_assistant/internal/reports"
)

const (
	tgMessageLimit            = 4000
	tgTextMaxLength           = 3500
	tgHTMLTagOpenHiddenQuote  = "<blockquote expandable>"
	tgHTMLTagCloseHiddenQuote = "</blockquote>"

	tgHTMLTagOpenQuote  = "<blockquote>"
	tgHTMLTagCloseQuote = "</blockquote>"

	tgHTMLTagOpenCode  = "<pre><code>"
	tgHTMLTagCloseCode = "</code></pre>"

	tgHTMLTagOpenBold  = "<b>"
	tgHTMLTagCloseBold = "</b>"
)

type TelegramBotRenderMode int

const (
	TelegramBotRenderHTML TelegramBotRenderMode = iota
	TelegramBotRenderPlain
)

type TelegramBotRenderer struct {
	Mode          TelegramBotRenderMode
	builder       *strings.Builder
	messages      []string
	isHiddenQuote bool
	isQuote       bool
	isCode        bool
	isBold        bool
}

func (r *TelegramBotRenderer) Render(data *reports.ReportData) ([]string, error) {
	r.builder = &strings.Builder{}
	r.builder.Grow(tgMessageLimit)
	r.messages = r.messages[:0]

	if data.Header != nil {
		if err := r.render(data.Header); err != nil {
			logger.Logf(logger.ERROR, "TelegramBotRenderer.Render()", "failed rendering header %v: %s", data.Header, err)
		} else {
			r.write("\n", false, true)
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

	if r.Mode == TelegramBotRenderHTML {
		return r.renderHTML(item)
	}

	r.renderPlain(item)
	return nil
}

func (r *TelegramBotRenderer) renderPlain(item *reports.Item) {
	for i, child := range item.Children {
		if i != 0 {
			if child.Block {
				r.write("\n", false, true)
			} else {
				r.write(" ", false, false)
			}
		}

		r.write(child.Text, false, false)

		if child.Link != "" {
			r.write("["+child.Link+"]", false, false)
		}

		if len(child.Children) > 0 {
			r.renderPlain(child)
		}
	}
}

func (r *TelegramBotRenderer) renderHTML(item *reports.Item) error {
	for i, child := range item.Children {
		if i != 0 {
			if child.Block {
				r.write("\n", false, true)
			} else {
				r.write(" ", false, false)
			}
		}

		r.tag(child.HiddenQuote, tgHTMLTagOpenHiddenQuote, &r.isHiddenQuote, true)
		r.tag(child.Quote, tgHTMLTagOpenQuote, &r.isQuote, true)
		r.tag(child.Code, tgHTMLTagOpenCode, &r.isCode, true)
		r.tag(child.Bold, tgHTMLTagOpenBold, &r.isBold, true)

		if child.Link != "" {
			r.write(`<a href="`+child.Link+`">`+child.Text+"</a>", false, false)
		} else {
			r.write(child.Text, false, false)
		}

		if len(child.Children) > 0 {
			if err := r.renderHTML(child); err != nil {
				return err
			}
		}

		r.tag(child.Bold, tgHTMLTagCloseBold, &r.isBold, false)
		r.tag(child.HiddenQuote, tgHTMLTagCloseHiddenQuote, &r.isHiddenQuote, false)
		r.tag(child.Quote, tgHTMLTagCloseQuote, &r.isQuote, false)
		r.tag(child.Code, tgHTMLTagCloseCode, &r.isCode, false)
	}
	return nil
}

func (r *TelegramBotRenderer) tag(flag bool, tag string, state *bool, value bool) {
	if flag {
		r.write(tag, true, false)
		*state = value
	}
}

func (r *TelegramBotRenderer) write(text string, tagOpen, tagClose bool) {
	b := r.builder
	builderLen := r.builder.Len()
	textLen := len(text)

	if textLen >= tgTextMaxLength {
		logger.Logf(logger.WARN, "TelegramBotRenderer.write()",
			"text length > %d, got %d, text: %s", tgTextMaxLength, textLen, text)
	}

	if tagOpen {
		if builderLen+textLen >= tgMessageLimit {
			r.closeTags()
			r.flush()
			r.openTags()
		}
		b.WriteString(text)
		return
	}

	if tagClose {
		b.WriteString(text)
		return
	}

	if builderLen+textLen >= tgMessageLimit {
		r.closeTags()
		r.flush()
		r.openTags()
	}

	b.WriteString(text)
}

func (r *TelegramBotRenderer) flush() {
	r.messages = append(r.messages, r.builder.String())
	r.builder.Reset()
}

func (r *TelegramBotRenderer) openTags() {
	if r.isHiddenQuote {
		r.write(tgHTMLTagOpenHiddenQuote, true, false)
	}
	if r.isQuote {
		r.write(tgHTMLTagOpenQuote, true, false)
	}
	if r.isCode {
		r.write(tgHTMLTagOpenCode, true, false)
	}
	if r.isBold {
		r.write(tgHTMLTagOpenBold, true, false)
	}
}

func (r *TelegramBotRenderer) closeTags() {
	if r.isBold {
		r.write(tgHTMLTagCloseBold, false, true)
	}
	if r.isHiddenQuote {
		r.write(tgHTMLTagCloseHiddenQuote, false, true)
	}
	if r.isQuote {
		r.write(tgHTMLTagCloseQuote, false, true)
	}
	if r.isCode {
		r.write(tgHTMLTagCloseCode, false, true)
	}
}
