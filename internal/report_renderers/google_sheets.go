package report_renderers

import (
	"wb_logistic_assistant/internal/reports"
)

const googleSheetsMaxWidth = 26

type GoogleSheetsRenderer struct {
	out  [][]interface{}
	posX int
	posY int
}

func (r *GoogleSheetsRenderer) Render(report *reports.ReportData) ([][]interface{}, error) {
	r.out = [][]interface{}{}
	r.posX, r.posY = 0, 0

	if report.Header != nil {
		r.render(report.Header)
		r.posY++
		r.posX = 0
	}

	if report.Body != nil {
		r.render(report.Body)
	}

	return r.out, nil
}

func (r *GoogleSheetsRenderer) render(item *reports.Item) {
	for i, child := range item.Children {
		if child == nil {
			r.expand(r.posY, r.posX)
			r.posX++
			continue
		}
		if child.Block && i != 0 {
			r.posY++
			r.posX = 0
		}

		if r.posX < googleSheetsMaxWidth && child.Text != "" || child.Link != "" {
			var val interface{}

			if child.Text != "" {
				val = child.Text
			}

			if child.Link != "" {
				if child.Text != "" {
					val = "=ГИПЕРССЫЛКА(\"" + child.Link + "\"; \"" + child.Text + "\")"
				} else {
					val = "=ГИПЕРССЫЛКА(\"" + child.Link + "\"; \"" + child.Link + "\")"
				}
			}

			r.expand(r.posY, r.posX)
			r.out[r.posY][r.posX] = val
			r.posX++
		}

		if len(child.Children) > 0 {
			r.render(child)
		}
	}
}

func (r *GoogleSheetsRenderer) expand(y, x int) {
	for len(r.out) <= y {
		r.out = append(r.out, make([]interface{}, 0))
	}

	row := r.out[y]
	if len(row) <= x {
		newRow := make([]interface{}, x+1)
		copy(newRow, row)
		r.out[y] = newRow
	}
}
