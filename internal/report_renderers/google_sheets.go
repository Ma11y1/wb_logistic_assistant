package report_renderers

import (
	"wb_logistic_assistant/internal/reports"
)

const googleSheetsMaxWidth = 26

type GoogleSheetsRenderer struct{}

func (*GoogleSheetsRenderer) Render(report *reports.ReportData) ([][]interface{}, error) {
	//headerLen := len(report.Header)
	//bodyLen := len(report.Body)
	//totalRows := headerLen + bodyLen
	//
	//output := make([][]interface{}, totalRows)
	//
	//copyRow := func(src []reports.Item) []interface{} {
	//	if src == nil {
	//		return nil
	//	}
	//	width := len(src)
	//	if width > googleSheetsMaxWidth {
	//		width = googleSheetsMaxWidth
	//	}
	//
	//	row := make([]interface{}, width)
	//	for i := 0; i < width; i++ {
	//		v := src[i]
	//		if v.Link != "" {
	//			row[i] = fmt.Sprintf(`=ГИПЕРССЫЛКА("%s","%s")`, v.Link, v.Text)
	//		} else {
	//			row[i] = v.Text
	//		}
	//	}
	//	return row
	//}
	//
	//for i := 0; i < headerLen; i++ {
	//	output[i] = copyRow(report.Header[i])
	//}
	//
	//for i := 0; i < bodyLen; i++ {
	//	output[headerLen+i] = copyRow(report.Body[i])
	//}
	//
	//return output, nil
	return nil, nil
}
