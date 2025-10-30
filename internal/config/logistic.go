package config

import (
	"encoding/json"
	"strconv"
	"time"
)

const logisticTimePeriod = time.Millisecond

type Logistic struct {
	wbClient *LogisticClient   // ro
	office   *LogisticOffice   // ro
	cacheTTL *LogisticCacheTTL // ro
}

type logistic struct {
	WBClient *LogisticClient   `json:"wb_client"`
	Office   *LogisticOffice   `json:"office"`
	CacheTTL *LogisticCacheTTL `json:"cache_ttl"`
}

func newLogistic() *Logistic {
	return &Logistic{
		wbClient: newLogisticClient(),   // default
		office:   newLogisticOffice(),   // default
		cacheTTL: newLogisticCacheTTL(), // default
	}
}

func (l *Logistic) WBClient() *LogisticClient   { return l.wbClient }
func (l *Logistic) Office() *LogisticOffice     { return l.office }
func (l *Logistic) CacheTTL() *LogisticCacheTTL { return l.cacheTTL }

func (l *Logistic) UnmarshalJSON(b []byte) error {
	temp := &logistic{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	l.wbClient = temp.WBClient
	l.office = temp.Office
	l.cacheTTL = temp.CacheTTL
	return nil
}

func (l *Logistic) MarshalJSON() ([]byte, error) {
	return json.Marshal(&logistic{
		WBClient: l.wbClient,
		Office:   l.office,
		CacheTTL: l.cacheTTL,
	})
}

type LogisticClient struct {
	userAgent    string // ro
	secUserAgent string // ro
	platform     string // ro
}

type logisticClient struct {
	UserAgent    string `json:"user_agent"`
	SecUserAgent string `json:"sec_user_agent"`
	Platform     string `json:"platform"`
}

func newLogisticClient() *LogisticClient {
	return &LogisticClient{
		userAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
		secUserAgent: "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"",
		platform:     "windows",
	}
}

func (l *LogisticClient) UserAgent() string    { return l.userAgent }
func (l *LogisticClient) SecUserAgent() string { return l.secUserAgent }
func (l *LogisticClient) Platform() string     { return l.platform }

func (l *LogisticClient) UnmarshalJSON(b []byte) error {
	temp := &logisticClient{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	l.userAgent = temp.UserAgent
	l.secUserAgent = temp.SecUserAgent
	l.platform = temp.Platform
	return nil
}

func (l *LogisticClient) MarshalJSON() ([]byte, error) {
	return json.Marshal(&logisticClient{
		UserAgent:    l.userAgent,
		SecUserAgent: l.secUserAgent,
		Platform:     l.platform,
	})
}

type LogisticOffice struct {
	id        int         // ro
	suppliers []int       // ro
	parking   map[int]int // ro
}

type logisticOffice struct {
	ID        int            `json:"id"`
	Suppliers []int          `json:"suppliers"`
	Parking   map[string]int `json:"parking"`
}

func newLogisticOffice() *LogisticOffice {
	return &LogisticOffice{
		id:        0,                 // default
		suppliers: []int{0},          // default
		parking:   make(map[int]int), // default
	}
}

func (l *LogisticOffice) UnmarshalJSON(b []byte) error {
	temp := &logisticOffice{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	l.id = temp.ID
	l.suppliers = temp.Suppliers
	l.parking = make(map[int]int)
	for key, parking := range temp.Parking {
		route, err := strconv.Atoi(key)
		if err != nil {
			return err
		}
		l.parking[route] = parking
	}
	return nil
}

func (l *LogisticOffice) MarshalJSON() ([]byte, error) {
	tempParking := make(map[string]int)
	for key, parking := range l.parking {
		tempParking[strconv.Itoa(key)] = parking
	}
	return json.Marshal(&logisticOffice{
		ID:        l.id,
		Suppliers: l.suppliers,
		Parking:   tempParking,
	})
}

func (l *LogisticOffice) ID() int              { return l.id }
func (l *LogisticOffice) Suppliers() []int     { return l.suppliers }
func (l *LogisticOffice) Parking() map[int]int { return l.parking }

type LogisticCacheTTL struct {
	userInfo                        time.Duration // ro
	remainsLastMileReports          time.Duration // ro
	remainsLastMileReportsRouteInfo time.Duration // ro
	jobsScheduling                  time.Duration // ro
	shipmentInfo                    time.Duration // ro
	shipmentTransfers               time.Duration // ro
	waySheetInfo                    time.Duration // ro
}

type logisticTTL struct {
	UserInfo                        time.Duration `json:"user_info"`
	RemainsLastMileReports          time.Duration `json:"remains_last_mile_report"`
	RemainsLastMileReportsRouteInfo time.Duration `json:"remains_last_mile_report_route_info"`
	JobsScheduling                  time.Duration `json:"jobs_scheduling"`
	ShipmentInfo                    time.Duration `json:"shipment_info"`
	ShipmentTransfers               time.Duration `json:"shipment_transfers"`
	WaySheetInfo                    time.Duration `json:"way_sheet_info"`
}

func newLogisticCacheTTL() *LogisticCacheTTL {
	return &LogisticCacheTTL{
		userInfo:                        3600000 * logisticTimePeriod, // default
		remainsLastMileReports:          60000 * logisticTimePeriod,   // default
		remainsLastMileReportsRouteInfo: 600000 * logisticTimePeriod,  // default
		jobsScheduling:                  1000000 * logisticTimePeriod, // default
		shipmentInfo:                    600000 * logisticTimePeriod,  // default
		shipmentTransfers:               120000 * logisticTimePeriod,  // default
		waySheetInfo:                    600000 * logisticTimePeriod,  // default
	}
}

func (l *LogisticCacheTTL) UserInfo() time.Duration               { return l.userInfo }
func (l *LogisticCacheTTL) RemainsLastMileReports() time.Duration { return l.remainsLastMileReports }
func (l *LogisticCacheTTL) RemainsLastMileReportsRouteInfo() time.Duration {
	return l.remainsLastMileReportsRouteInfo
}
func (l *LogisticCacheTTL) JobsScheduling() time.Duration    { return l.jobsScheduling }
func (l *LogisticCacheTTL) ShipmentInfo() time.Duration      { return l.shipmentInfo }
func (l *LogisticCacheTTL) ShipmentTransfers() time.Duration { return l.shipmentTransfers }
func (l *LogisticCacheTTL) WaySheetInfo() time.Duration      { return l.waySheetInfo }

func (l *LogisticCacheTTL) UnmarshalJSON(b []byte) error {
	temp := &logisticTTL{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	l.userInfo = temp.UserInfo * logisticTimePeriod
	l.remainsLastMileReports = temp.RemainsLastMileReports * logisticTimePeriod
	l.remainsLastMileReportsRouteInfo = temp.RemainsLastMileReportsRouteInfo * logisticTimePeriod
	l.jobsScheduling = temp.JobsScheduling * logisticTimePeriod
	l.shipmentInfo = temp.ShipmentInfo * logisticTimePeriod
	l.shipmentTransfers = temp.ShipmentTransfers * logisticTimePeriod
	l.waySheetInfo = temp.WaySheetInfo * logisticTimePeriod
	return nil
}

func (l *LogisticCacheTTL) MarshalJSON() ([]byte, error) {
	return json.Marshal(&logisticTTL{
		UserInfo:                        l.userInfo / logisticTimePeriod,
		RemainsLastMileReports:          l.remainsLastMileReports / logisticTimePeriod,
		RemainsLastMileReportsRouteInfo: l.remainsLastMileReportsRouteInfo / logisticTimePeriod,
		JobsScheduling:                  l.jobsScheduling / logisticTimePeriod,
		ShipmentInfo:                    l.shipmentInfo / logisticTimePeriod,
		ShipmentTransfers:               l.shipmentTransfers / logisticTimePeriod,
		WaySheetInfo:                    l.waySheetInfo / logisticTimePeriod,
	})
}
