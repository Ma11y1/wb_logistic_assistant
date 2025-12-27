package report_renderers

import "wb_logistic_assistant/internal/reports"

type ReportRenderer[T interface{}] interface {
	Render(data *reports.ReportData) (T, error)
}
