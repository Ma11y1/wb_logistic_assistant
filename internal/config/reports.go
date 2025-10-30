package config

import (
	"encoding/json"
	"time"
)

const reportsTimePeriod = time.Millisecond

type Reports struct {
	generalRoutes *ReportsGeneralRoutes // ro
	shipmentClose *ReportsShipmentClose // ro
}

type reports struct {
	GeneralRoutes *ReportsGeneralRoutes `json:"general_routes"`
	ShipmentClose *ReportsShipmentClose `json:"shipment_close"`
}

func newReports() *Reports {
	return &Reports{
		generalRoutes: newReportsGeneralRoutes(), // default
		shipmentClose: newReportsShipmentClose(), // default
	}
}

func (r *Reports) GeneralRoutes() *ReportsGeneralRoutes { return r.generalRoutes }
func (r *Reports) ShipmentClose() *ReportsShipmentClose { return r.shipmentClose }

func (r *Reports) UnmarshalJSON(b []byte) error {
	temp := &reports{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	r.generalRoutes = temp.GeneralRoutes
	r.shipmentClose = temp.ShipmentClose
	return nil
}

func (r *Reports) MarshalJSON() ([]byte, error) {
	return json.Marshal(&reports{
		GeneralRoutes: r.generalRoutes,
		ShipmentClose: r.shipmentClose,
	})
}

type ReportsGeneralRoutes struct {
	isEnabled                    bool          // ro
	pollingInterval              time.Duration // ro
	errRetryTaskLimit            int           // ro
	taskTimeout                  time.Duration // ro
	ttlResetChangeBarcodes       time.Duration // ro
	ttlLoadRemainsLastMileReport time.Duration // ro
	ttlLoadJobsScheduling        time.Duration // ro
	ttlLoadShipments             time.Duration // ro
	ttlLoadWaySheets             time.Duration // to
	isSort                       bool          // ro
	isSortAscending              bool          // ro
	sortColumn                   int           // ro
	skipRoutes                   []int         // ro
	isRenderGoogleSheets         bool          // ro
	isRenderTelegramBot          bool          // ro
}

type reportsGeneralRoutes struct {
	IsEnabled                    bool          `json:"enabled"`
	PollingInterval              time.Duration `json:"polling_interval"`
	TaskTimeout                  time.Duration `json:"task_timeout"`
	ErrRetryTaskLimit            int           `json:"err_retry_task_limit"`
	TTLResetChangeBarcodes       time.Duration `json:"ttl_reset_change_barcodes"`
	TTLLoadRemainsLastMileReport time.Duration `json:"ttl_load_remains_last_mile_report"`
	TTLLoadJobsScheduling        time.Duration `json:"ttl_load_jobs_scheduling"`
	TTLLoadShipments             time.Duration `json:"ttl_load_shipments"`
	TTLLoadWaySheets             time.Duration `json:"ttl_load_way_sheets"`
	IsSort                       bool          `json:"sort"`
	IsSortAscending              bool          `json:"sort_ascending"`
	SortColumn                   int           `json:"sort_column"`
	SkipRoutes                   []int         `json:"skip_routes"`
	IsRenderGoogleSheets         bool          `json:"render_google_sheets"`
	IsRenderTelegramBot          bool          `json:"render_telegram_bot"`
}

func newReportsGeneralRoutes() *ReportsGeneralRoutes {
	return &ReportsGeneralRoutes{
		isEnabled:                    false,                       // default
		pollingInterval:              5000 * reportsTimePeriod,    // default
		taskTimeout:                  600000 * reportsTimePeriod,  // default
		errRetryTaskLimit:            3,                           // default
		ttlResetChangeBarcodes:       60000 * reportsTimePeriod,   // default
		ttlLoadRemainsLastMileReport: 20000 * reportsTimePeriod,   // default
		ttlLoadJobsScheduling:        1000000 * reportsTimePeriod, // default
		ttlLoadShipments:             300000 * reportsTimePeriod,  // default
		ttlLoadWaySheets:             180000 * reportsTimePeriod,  // default
		isSort:                       false,                       // default
		isSortAscending:              false,                       // default
		sortColumn:                   0,                           // default
		skipRoutes:                   []int{},                     // default
		isRenderGoogleSheets:         false,                       // default
		isRenderTelegramBot:          false,                       // default
	}
}

func (t *ReportsGeneralRoutes) IsEnabled() bool { return t.isEnabled }

func (t *ReportsGeneralRoutes) PollingInterval() time.Duration { return t.pollingInterval }

func (t *ReportsGeneralRoutes) TaskTimeout() time.Duration { return t.taskTimeout }

func (t *ReportsGeneralRoutes) ErrRetryTaskLimit() int { return t.errRetryTaskLimit }

func (t *ReportsGeneralRoutes) TTLResetChangeBarcodes() time.Duration {
	return t.ttlResetChangeBarcodes
}
func (t *ReportsGeneralRoutes) TTLLoadRemainsLastMileReport() time.Duration {
	return t.ttlLoadRemainsLastMileReport
}
func (t *ReportsGeneralRoutes) TTLLoadJobsScheduling() time.Duration {
	return t.ttlLoadJobsScheduling
}
func (t *ReportsGeneralRoutes) TTLLoadShipments() time.Duration {
	return t.ttlLoadShipments
}
func (t *ReportsGeneralRoutes) TTLLoadWaySheets() time.Duration {
	return t.ttlLoadWaySheets
}

func (t *ReportsGeneralRoutes) IsSort() bool      { return t.isSort }
func (t *ReportsGeneralRoutes) IsAscending() bool { return t.isSortAscending }
func (t *ReportsGeneralRoutes) SortColumn() int   { return t.sortColumn }

func (t *ReportsGeneralRoutes) SkipRoutes() []int { return t.skipRoutes }

func (t *ReportsGeneralRoutes) IsRenderGoogleSheets() bool { return t.isRenderGoogleSheets }
func (t *ReportsGeneralRoutes) IsRenderTelegramBot() bool  { return t.isRenderTelegramBot }

func (t *ReportsGeneralRoutes) UnmarshalJSON(b []byte) error {
	temp := &reportsGeneralRoutes{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	t.isEnabled = temp.IsEnabled
	t.pollingInterval = temp.PollingInterval * reportsTimePeriod
	t.errRetryTaskLimit = temp.ErrRetryTaskLimit
	t.taskTimeout = temp.TaskTimeout * reportsTimePeriod
	t.ttlResetChangeBarcodes = temp.TTLResetChangeBarcodes * reportsTimePeriod
	t.ttlLoadRemainsLastMileReport = temp.TTLLoadRemainsLastMileReport * reportsTimePeriod
	t.ttlLoadJobsScheduling = temp.TTLLoadJobsScheduling * reportsTimePeriod
	t.ttlLoadShipments = temp.TTLLoadShipments * reportsTimePeriod
	t.isSort = temp.IsSort
	t.isSortAscending = temp.IsSortAscending
	t.sortColumn = temp.SortColumn
	t.skipRoutes = temp.SkipRoutes
	t.isRenderGoogleSheets = temp.IsRenderGoogleSheets
	t.isRenderTelegramBot = temp.IsRenderTelegramBot
	return nil
}

func (t *ReportsGeneralRoutes) MarshalJSON() ([]byte, error) {
	return json.Marshal(&reportsGeneralRoutes{
		IsEnabled:                    t.isEnabled,
		PollingInterval:              t.pollingInterval / reportsTimePeriod,
		ErrRetryTaskLimit:            t.errRetryTaskLimit,
		TaskTimeout:                  t.taskTimeout / reportsTimePeriod,
		TTLResetChangeBarcodes:       t.ttlResetChangeBarcodes / reportsTimePeriod,
		TTLLoadRemainsLastMileReport: t.ttlLoadRemainsLastMileReport / reportsTimePeriod,
		TTLLoadJobsScheduling:        t.ttlLoadJobsScheduling / reportsTimePeriod,
		TTLLoadShipments:             t.ttlLoadShipments / reportsTimePeriod,
		IsSort:                       t.isSort,
		IsSortAscending:              t.isSortAscending,
		SortColumn:                   t.sortColumn,
		SkipRoutes:                   t.skipRoutes,
		IsRenderGoogleSheets:         t.isRenderGoogleSheets,
		IsRenderTelegramBot:          t.isRenderTelegramBot,
	})
}

type ReportsShipmentClose struct {
	isEnabled             bool          // ro
	pollingInterval       time.Duration // ro
	findShipmentsInterval time.Duration // ro
	errRetryTaskLimit     int           // ro
	taskTimeout           time.Duration // ro
	skipRoutes            []int         // ro
	isRenderGoogleSheets  bool          // ro
	isRenderTelegramBot   bool          // ro
}

type reportsShipmentClose struct {
	IsEnabled             bool          `json:"enabled"`
	PollingInterval       time.Duration `json:"polling_interval"`
	FindShipmentsInterval time.Duration `json:"find_shipments_interval"`
	TaskTimeout           time.Duration `json:"task_timeout"`
	ErrRetryTaskLimit     int           `json:"err_retry_task_limit"`
	SkipRoutes            []int         `json:"skip_routes"`
	IsRenderGoogleSheets  bool          `json:"render_google_sheets"`
	IsRenderTelegramBot   bool          `json:"render_telegram_bot"`
}

func newReportsShipmentClose() *ReportsShipmentClose {
	return &ReportsShipmentClose{
		isEnabled:             false,                      // default
		pollingInterval:       1000 * reportsTimePeriod,   // default
		findShipmentsInterval: 10000 * reportsTimePeriod,  // default
		taskTimeout:           600000 * reportsTimePeriod, // default
		errRetryTaskLimit:     3,                          // default
		skipRoutes:            []int{},                    // default
		isRenderGoogleSheets:  false,                      // default
		isRenderTelegramBot:   false,                      // default
	}
}

func (t *ReportsShipmentClose) IsEnabled() bool { return t.isEnabled }

func (t *ReportsShipmentClose) PollingInterval() time.Duration { return t.pollingInterval }

func (t *ReportsShipmentClose) FindShipmentsInterval() time.Duration { return t.findShipmentsInterval }

func (t *ReportsShipmentClose) TaskTimeout() time.Duration { return t.taskTimeout }

func (t *ReportsShipmentClose) ErrRetryTaskLimit() int { return t.errRetryTaskLimit }

func (t *ReportsShipmentClose) SkipRoutes() []int { return t.skipRoutes }

func (t *ReportsShipmentClose) IsRenderGoogleSheets() bool { return t.isRenderGoogleSheets }
func (t *ReportsShipmentClose) IsRenderTelegramBot() bool  { return t.isRenderTelegramBot }

func (t *ReportsShipmentClose) UnmarshalJSON(b []byte) error {
	temp := &reportsShipmentClose{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	t.isEnabled = temp.IsEnabled
	t.pollingInterval = temp.PollingInterval * reportsTimePeriod
	t.findShipmentsInterval = temp.FindShipmentsInterval * reportsTimePeriod
	t.errRetryTaskLimit = temp.ErrRetryTaskLimit
	t.taskTimeout = temp.TaskTimeout * reportsTimePeriod
	t.skipRoutes = temp.SkipRoutes
	t.isRenderGoogleSheets = temp.IsRenderGoogleSheets
	t.isRenderTelegramBot = temp.IsRenderTelegramBot
	return nil
}

func (t *ReportsShipmentClose) MarshalJSON() ([]byte, error) {
	return json.Marshal(&reportsShipmentClose{
		IsEnabled:             t.isEnabled,
		PollingInterval:       t.pollingInterval / reportsTimePeriod,
		FindShipmentsInterval: t.findShipmentsInterval / reportsTimePeriod,
		ErrRetryTaskLimit:     t.errRetryTaskLimit,
		TaskTimeout:           t.taskTimeout / reportsTimePeriod,
		SkipRoutes:            t.skipRoutes,
		IsRenderGoogleSheets:  t.isRenderGoogleSheets,
		IsRenderTelegramBot:   t.isRenderTelegramBot,
	})
}
