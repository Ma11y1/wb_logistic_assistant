package models

type CaptchaBrowserInfo struct {
	UserAgent struct {
		Value string `json:"value"`
		State string `json:"state"`
	} `json:"userAgent"`
	AppVersion struct {
		Value string `json:"value"`
		State string `json:"state"`
	} `json:"appVersion"`
	Languages struct {
		Value struct {
			Main                          string   `json:"main"`
			Preferred                     []string `json:"preferred"`
			LocaleLanguagesDateTimeFormat string   `json:"localeLanguagesDateTimeFormat"`
		} `json:"value"`
		State string `json:"state"`
	} `json:"languages"`
	MimeTypes struct {
		Value struct {
			MimeTypes              []string `json:"mimeTypes"`
			IsPrototypesConsistent bool     `json:"isPrototypesConsistent"`
		} `json:"value"`
		State string `json:"state"`
	} `json:"mimeTypes"`
	EvalLength struct {
		Value int    `json:"value"`
		State string `json:"state"`
	} `json:"evalLength"`
	Plugins struct {
		Value struct {
			Plugins             []string `json:"plugins"`
			PluginsTrueInstance bool     `json:"pluginsTrueInstance"`
		} `json:"value"`
		State string `json:"state"`
	} `json:"plugins"`
	DocumentElementAttrs struct {
		Value []string `json:"value"`
		State string   `json:"state"`
	} `json:"documentElementAttrs"`
	ErrorTrace struct {
		Value string `json:"value"`
		State string `json:"state"`
	} `json:"errorTrace"`
	FunctionBind struct {
		Value string `json:"value"`
		State string `json:"state"`
	} `json:"functionBind"`
	ProductSub struct {
		Value string `json:"value"`
		State string `json:"state"`
	} `json:"productSub"`
	WebDriver struct {
		Value bool   `json:"value"`
		State string `json:"state"`
	} `json:"webDriver"`
	WebGL struct {
		Value struct {
			Vendor   string `json:"vendor"`
			Renderer string `json:"renderer"`
		} `json:"value"`
		State string `json:"state"`
	} `json:"webGL"`
	Screen struct {
		Value struct {
			WInnerHeight      int `json:"wInnerHeight"`
			WOuterHeight      int `json:"wOuterHeight"`
			WOuterWidth       int `json:"wOuterWidth"`
			WInnerWidth       int `json:"wInnerWidth"`
			WScreenX          int `json:"wScreenX"`
			WPageXOffset      int `json:"wPageXOffset"`
			WPageYOffset      int `json:"wPageYOffset"`
			CWidth            int `json:"cWidth"`
			CHeight           int `json:"cHeight"`
			SWidth            int `json:"sWidth"`
			SHeight           int `json:"sHeight"`
			SAvailWidth       int `json:"sAvailWidth"`
			SAvailHeight      int `json:"sAvailHeight"`
			SColorDepth       int `json:"sColorDepth"`
			SPixelDepth       int `json:"sPixelDepth"`
			WDevicePixelRatio int `json:"wDevicePixelRatio"`
		} `json:"value"`
		State string `json:"state"`
	} `json:"screen"`
	BrowserVendor struct {
		Value string `json:"value"`
		State string `json:"state"`
	} `json:"browserVendor"`
	BrowserEngineType struct {
		Value string `json:"value"`
		State string `json:"state"`
	} `json:"browserEngineType"`
	Permissions struct {
		Value struct {
			InconsistentNotifications bool              `json:"inconsistentNotifications"`
			Permissions               map[string]string `json:"permissions"`
		} `json:"value"`
		State string `json:"state"`
	} `json:"permissions"`
}

type CaptchaTask struct {
	Timestamp    int    `json:"timestamp"`
	Value        string `json:"value"`
	ThresholdHex string `json:"threshold_hex"`
	MatchBits    int    `json:"match_bits"`
	BufferLen    int    `json:"buffer_len"`
	AnswersCount int    `json:"answers_count"`
	R            int    `json:"r"`
	N            int    `json:"n"`
	ClientId     string `json:"client_id"`
	Sign         string `json:"sign"`
}
