package reporters

type SortType int

const (
	IDSortType SortType = iota
	ParkingSortType
	TaresSortType
	BarcodesSortType
	RatingSortType
)
