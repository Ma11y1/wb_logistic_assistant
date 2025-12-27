package reports

type ReportData struct {
	Header *Item `json:"header"`
	Body   *Item `json:"body"`
}

func NewReportData() *ReportData {
	return &ReportData{
		Header: &Item{},
		Body:   &Item{},
	}
}

type Item struct {
	Text        string
	Link        string
	Bold        bool // style, telegram
	Quote       bool // block, telegram
	HiddenQuote bool // block, telegram
	Code        bool // block, telegram
	Block       bool
	Children    []*Item
}

func NewItem() *Item {
	return &Item{Children: []*Item{}}
}

func (i *Item) AddChild(child *Item) {
	if child == nil {
		return
	}
	i.Children = append(i.Children, child)
}

func (i *Item) RemoveChild(child *Item) {
	if child == nil {
		return
	}
	for j, c := range i.Children {
		if c == child {
			i.Children = append(i.Children[:j], i.Children[j+1:]...)
		}
	}
}

func (i *Item) CountChild() int {
	return len(i.Children) - 1
}
