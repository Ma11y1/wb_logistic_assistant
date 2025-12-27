package config

import (
	"encoding/json"
	"time"
)

const reportsTimePeriod = time.Millisecond

type Reports struct {
	generalRoutes *ReportsGeneralRoutes // ro
	shipmentClose *ReportsShipmentClose // ro
	financeRoutes *ReportsFinanceRoutes // ro
	financeDaily  *ReportsFinanceDaily  // ro
}

type reports struct {
	GeneralRoutes *ReportsGeneralRoutes `json:"general_routes"`
	ShipmentClose *ReportsShipmentClose `json:"shipment_close"`
	FinanceRoutes *ReportsFinanceRoutes `json:"finance_routes"`
	FinanceDaily  *ReportsFinanceDaily  `json:"finance_daily"`
}

func newReports() *Reports {
	return &Reports{
		generalRoutes: newReportsGeneralRoutes(), // default
		shipmentClose: newReportsShipmentClose(), // default
		financeRoutes: newReportsFinanceRoutes(), // default
		financeDaily:  newReportsFinanceDaily(),  // default
	}
}

func (r *Reports) GeneralRoutes() *ReportsGeneralRoutes { return r.generalRoutes }
func (r *Reports) ShipmentClose() *ReportsShipmentClose { return r.shipmentClose }
func (r *Reports) FinanceRoutes() *ReportsFinanceRoutes { return r.financeRoutes }
func (r *Reports) FinanceDaily() *ReportsFinanceDaily   { return r.financeDaily }

func (r *Reports) UnmarshalJSON(b []byte) error {
	temp := &reports{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	r.generalRoutes = temp.GeneralRoutes
	r.shipmentClose = temp.ShipmentClose
	r.financeRoutes = temp.FinanceRoutes
	r.financeDaily = temp.FinanceDaily
	return nil
}

func (r *Reports) MarshalJSON() ([]byte, error) {
	return json.Marshal(&reports{
		GeneralRoutes: r.generalRoutes,
		ShipmentClose: r.shipmentClose,
		FinanceRoutes: r.financeRoutes,
		FinanceDaily:  r.financeDaily,
	})
}

type ReportsGeneralRoutes struct {
	isEnabled                   bool          // ro
	errRetryTaskLimit           int           // ro
	pollingInterval             time.Duration // ro
	taskTimeout                 time.Duration // ro
	intervalResetChangeBarcodes time.Duration // ro
	intervalUpdateRating        time.Duration // ro
	intervalUpdateShipments     time.Duration // ro
	intervalUpdateWaySheets     time.Duration // ro
	isSort                      bool          // ro
	isSortAscending             bool          // ro
	sortColumn                  int           // ro
	isRenderGoogleSheets        bool          // ro
}

type reportsGeneralRoutes struct {
	IsEnabled                   bool          `json:"enabled"`
	ErrRetryTaskLimit           int           `json:"err_retry_task_limit"`
	PollingInterval             time.Duration `json:"polling_interval"`
	TaskTimeout                 time.Duration `json:"task_timeout"`
	IntervalResetChangeBarcodes time.Duration `json:"interval_reset_change_barcodes"`
	IntervalUpdateRating        time.Duration `json:"interval_update_rating"`
	IntervalUpdateShipments     time.Duration `json:"interval_update_shipments"`
	IntervalUpdateWaySheets     time.Duration `json:"interval_update_waysheets"`
	IsSort                      bool          `json:"sort"`
	IsSortAscending             bool          `json:"sort_ascending"`
	SortColumn                  int           `json:"sort_column"`
	IsRenderGoogleSheets        bool          `json:"render_google_sheets"`
}

func newReportsGeneralRoutes() *ReportsGeneralRoutes {
	return &ReportsGeneralRoutes{
		isEnabled:                   false,                         // default
		errRetryTaskLimit:           3,                             // default
		pollingInterval:             5000 * reportsTimePeriod,      // default
		taskTimeout:                 600_000 * reportsTimePeriod,   // default
		intervalResetChangeBarcodes: 300_000 * reportsTimePeriod,   // default
		intervalUpdateRating:        6_000_000 * reportsTimePeriod, // default
		intervalUpdateShipments:     100_000 * reportsTimePeriod,   // default
		intervalUpdateWaySheets:     150_000 * reportsTimePeriod,   // default
		isSort:                      false,                         // default
		isSortAscending:             false,                         // default
		sortColumn:                  0,                             // default
		isRenderGoogleSheets:        false,                         // default
	}
}

func (r *ReportsGeneralRoutes) IsEnabled() bool { return r.isEnabled }

func (r *ReportsGeneralRoutes) ErrRetryTaskLimit() int { return r.errRetryTaskLimit }

func (r *ReportsGeneralRoutes) PollingInterval() time.Duration { return r.pollingInterval }

func (r *ReportsGeneralRoutes) TaskTimeout() time.Duration { return r.taskTimeout }

func (r *ReportsGeneralRoutes) IntervalResetChangeBarcodes() time.Duration {
	return r.intervalResetChangeBarcodes
}
func (r *ReportsGeneralRoutes) IntervalUpdateRating() time.Duration {
	return r.intervalUpdateRating
}
func (r *ReportsGeneralRoutes) IntervalUpdateShipments() time.Duration {
	return r.intervalUpdateShipments
}
func (r *ReportsGeneralRoutes) IntervalUpdateWaySheets() time.Duration {
	return r.intervalUpdateWaySheets
}

func (r *ReportsGeneralRoutes) IsSort() bool      { return r.isSort }
func (r *ReportsGeneralRoutes) IsAscending() bool { return r.isSortAscending }
func (r *ReportsGeneralRoutes) SortColumn() int   { return r.sortColumn }

func (r *ReportsGeneralRoutes) IsRenderGoogleSheets() bool { return r.isRenderGoogleSheets }

func (r *ReportsGeneralRoutes) UnmarshalJSON(b []byte) error {
	temp := &reportsGeneralRoutes{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	r.isEnabled = temp.IsEnabled
	r.pollingInterval = temp.PollingInterval * reportsTimePeriod
	r.errRetryTaskLimit = temp.ErrRetryTaskLimit
	r.taskTimeout = temp.TaskTimeout * reportsTimePeriod
	r.intervalResetChangeBarcodes = temp.IntervalResetChangeBarcodes * reportsTimePeriod
	r.intervalUpdateRating = temp.IntervalUpdateRating * reportsTimePeriod
	r.intervalUpdateShipments = temp.IntervalUpdateShipments * reportsTimePeriod
	r.intervalUpdateWaySheets = temp.IntervalUpdateWaySheets * reportsTimePeriod
	r.isSort = temp.IsSort
	r.isSortAscending = temp.IsSortAscending
	r.sortColumn = temp.SortColumn
	r.isRenderGoogleSheets = temp.IsRenderGoogleSheets
	return nil
}

func (r *ReportsGeneralRoutes) MarshalJSON() ([]byte, error) {
	return json.Marshal(&reportsGeneralRoutes{
		IsEnabled:                   r.isEnabled,
		PollingInterval:             r.pollingInterval / reportsTimePeriod,
		ErrRetryTaskLimit:           r.errRetryTaskLimit,
		TaskTimeout:                 r.taskTimeout / reportsTimePeriod,
		IntervalResetChangeBarcodes: r.intervalResetChangeBarcodes / reportsTimePeriod,
		IntervalUpdateRating:        r.intervalUpdateRating / reportsTimePeriod,
		IntervalUpdateShipments:     r.intervalUpdateShipments / reportsTimePeriod,
		IntervalUpdateWaySheets:     r.intervalUpdateWaySheets / reportsTimePeriod,
		IsSort:                      r.isSort,
		IsSortAscending:             r.isSortAscending,
		SortColumn:                  r.sortColumn,
		IsRenderGoogleSheets:        r.isRenderGoogleSheets,
	})
}

type ReportsShipmentClose struct {
	isEnabled               bool          // ro
	errRetryTaskLimit       int           // ro
	pollingInterval         time.Duration // ro
	taskTimeout             time.Duration // ro
	intervalUpdateShipments time.Duration // ro
	isRenderGoogleSheets    bool          // ro
	isRenderTelegramBot     bool          // ro
}

type reportsShipmentClose struct {
	IsEnabled               bool          `json:"enabled"`
	ErrRetryTaskLimit       int           `json:"err_retry_task_limit"`
	PollingInterval         time.Duration `json:"polling_interval"`
	TaskTimeout             time.Duration `json:"task_timeout"`
	IntervalUpdateShipments time.Duration `json:"interval_update_shipments"`
	IsRenderGoogleSheets    bool          `json:"render_google_sheets"`
	IsRenderTelegramBot     bool          `json:"render_telegram_bot"`
}

func newReportsShipmentClose() *ReportsShipmentClose {
	return &ReportsShipmentClose{
		isEnabled:               false,                       // default
		errRetryTaskLimit:       3,                           // default
		pollingInterval:         1000 * reportsTimePeriod,    // default
		taskTimeout:             600_000 * reportsTimePeriod, // default
		intervalUpdateShipments: 100_000 * reportsTimePeriod, // default
		isRenderGoogleSheets:    false,                       // default
		isRenderTelegramBot:     false,                       // default
	}
}

func (r *ReportsShipmentClose) IsEnabled() bool { return r.isEnabled }

func (r *ReportsShipmentClose) PollingInterval() time.Duration { return r.pollingInterval }

func (r *ReportsShipmentClose) IntervalUpdateShipments() time.Duration {
	return r.intervalUpdateShipments
}

func (r *ReportsShipmentClose) TaskTimeout() time.Duration { return r.taskTimeout }

func (r *ReportsShipmentClose) ErrRetryTaskLimit() int { return r.errRetryTaskLimit }

func (r *ReportsShipmentClose) IsRenderGoogleSheets() bool { return r.isRenderGoogleSheets }
func (r *ReportsShipmentClose) IsRenderTelegramBot() bool  { return r.isRenderTelegramBot }

func (r *ReportsShipmentClose) UnmarshalJSON(b []byte) error {
	temp := &reportsShipmentClose{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	r.isEnabled = temp.IsEnabled
	r.pollingInterval = temp.PollingInterval * reportsTimePeriod
	r.intervalUpdateShipments = temp.IntervalUpdateShipments * reportsTimePeriod
	r.errRetryTaskLimit = temp.ErrRetryTaskLimit
	r.taskTimeout = temp.TaskTimeout * reportsTimePeriod
	r.isRenderGoogleSheets = temp.IsRenderGoogleSheets
	r.isRenderTelegramBot = temp.IsRenderTelegramBot
	return nil
}

func (r *ReportsShipmentClose) MarshalJSON() ([]byte, error) {
	return json.Marshal(&reportsShipmentClose{
		IsEnabled:               r.isEnabled,
		PollingInterval:         r.pollingInterval / reportsTimePeriod,
		IntervalUpdateShipments: r.intervalUpdateShipments / reportsTimePeriod,
		ErrRetryTaskLimit:       r.errRetryTaskLimit,
		TaskTimeout:             r.taskTimeout / reportsTimePeriod,
		IsRenderGoogleSheets:    r.isRenderGoogleSheets,
		IsRenderTelegramBot:     r.isRenderTelegramBot,
	})
}

type ReportsFinanceRoutes struct {
	isEnabled           bool          // ro
	errRetryTaskLimit   int           // ro
	pollingInterval     time.Duration // ro
	taskTimeout         time.Duration // ro
	isRenderTelegramBot bool          // ro
}

type reportsFinanceRoutes struct {
	IsEnabled           bool          `json:"enabled"`
	ErrRetryTaskLimit   int           `json:"err_retry_task_limit"`
	PollingInterval     time.Duration `json:"polling_interval"`
	TaskTimeout         time.Duration `json:"task_timeout"`
	IsRenderTelegramBot bool          `json:"render_telegram_bot"`
}

func newReportsFinanceRoutes() *ReportsFinanceRoutes {
	return &ReportsFinanceRoutes{
		isEnabled:           false,                       // default
		errRetryTaskLimit:   3,                           // default
		pollingInterval:     1000 * reportsTimePeriod,    // default
		taskTimeout:         600_000 * reportsTimePeriod, // default
		isRenderTelegramBot: false,                       // default
	}
}

func (r *ReportsFinanceRoutes) IsEnabled() bool { return r.isEnabled }

func (r *ReportsFinanceRoutes) PollingInterval() time.Duration { return r.pollingInterval }

func (r *ReportsFinanceRoutes) TaskTimeout() time.Duration { return r.taskTimeout }

func (r *ReportsFinanceRoutes) ErrRetryTaskLimit() int { return r.errRetryTaskLimit }

func (r *ReportsFinanceRoutes) IsRenderTelegramBot() bool { return r.isRenderTelegramBot }

func (r *ReportsFinanceRoutes) UnmarshalJSON(b []byte) error {
	temp := &reportsFinanceRoutes{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	r.isEnabled = temp.IsEnabled
	r.pollingInterval = temp.PollingInterval * reportsTimePeriod
	r.errRetryTaskLimit = temp.ErrRetryTaskLimit
	r.taskTimeout = temp.TaskTimeout * reportsTimePeriod
	r.isRenderTelegramBot = temp.IsRenderTelegramBot
	return nil
}

func (r *ReportsFinanceRoutes) MarshalJSON() ([]byte, error) {
	return json.Marshal(&reportsFinanceRoutes{
		IsEnabled:           r.isEnabled,
		PollingInterval:     r.pollingInterval / reportsTimePeriod,
		ErrRetryTaskLimit:   r.errRetryTaskLimit,
		TaskTimeout:         r.taskTimeout / reportsTimePeriod,
		IsRenderTelegramBot: r.isRenderTelegramBot,
	})
}

type ReportsFinanceDaily struct {
	isEnabled           bool          // ro
	errRetryTaskLimit   int           // ro
	pollingInterval     time.Duration // ro
	taskTimeout         time.Duration // ro
	renderAtStart       bool          // ro
	dayOffset           int           // ro
	isRenderTelegramBot bool          // ro
}

type reportsFinanceDaily struct {
	IsEnabled           bool          `json:"enabled"`
	ErrRetryTaskLimit   int           `json:"err_retry_task_limit"`
	PollingInterval     time.Duration `json:"polling_interval"`
	TaskTimeout         time.Duration `json:"task_timeout"`
	RenderAtStart       bool          `json:"render_at_start"`
	DayOffset           int           `json:"day_offset"`
	IsRenderTelegramBot bool          `json:"render_telegram_bot"`
}

func newReportsFinanceDaily() *ReportsFinanceDaily {
	return &ReportsFinanceDaily{
		isEnabled:           false,                       // default
		errRetryTaskLimit:   3,                           // default
		pollingInterval:     1000 * reportsTimePeriod,    // default
		taskTimeout:         600_000 * reportsTimePeriod, // default
		renderAtStart:       false,                       // default
		dayOffset:           -1,                          // default
		isRenderTelegramBot: false,                       // default
	}
}

func (r *ReportsFinanceDaily) IsEnabled() bool { return r.isEnabled }

func (r *ReportsFinanceDaily) PollingInterval() time.Duration { return r.pollingInterval }

func (r *ReportsFinanceDaily) TaskTimeout() time.Duration { return r.taskTimeout }

func (r *ReportsFinanceDaily) ErrRetryTaskLimit() int { return r.errRetryTaskLimit }

func (r *ReportsFinanceDaily) RenderAtStart() bool {
	return r.renderAtStart
}

func (r *ReportsFinanceDaily) DayOffset() int { return r.dayOffset }

func (r *ReportsFinanceDaily) IsRenderTelegramBot() bool { return r.isRenderTelegramBot }

func (r *ReportsFinanceDaily) UnmarshalJSON(b []byte) error {
	temp := &reportsFinanceDaily{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	r.isEnabled = temp.IsEnabled
	r.pollingInterval = temp.PollingInterval * reportsTimePeriod
	r.errRetryTaskLimit = temp.ErrRetryTaskLimit
	r.taskTimeout = temp.TaskTimeout * reportsTimePeriod
	r.renderAtStart = temp.RenderAtStart
	r.dayOffset = temp.DayOffset
	r.isRenderTelegramBot = temp.IsRenderTelegramBot
	return nil
}

func (r *ReportsFinanceDaily) MarshalJSON() ([]byte, error) {
	return json.Marshal(&reportsFinanceDaily{
		IsEnabled:           r.isEnabled,
		PollingInterval:     r.pollingInterval / reportsTimePeriod,
		ErrRetryTaskLimit:   r.errRetryTaskLimit,
		TaskTimeout:         r.taskTimeout / reportsTimePeriod,
		RenderAtStart:       r.renderAtStart,
		DayOffset:           r.dayOffset,
		IsRenderTelegramBot: r.isRenderTelegramBot,
	})
}
