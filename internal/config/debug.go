package config

import "encoding/json"

type Debug struct {
	enabled bool
	path    string
}

type debug struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

func newDebug() *Debug {
	return &Debug{
		enabled: false,
		path:    "debug.txt",
	}
}

func (d *Debug) Enabled() bool { return d.enabled }
func (d *Debug) Path() string  { return d.path }

func (d *Debug) UnmarshalJSON(b []byte) error {
	temp := &debug{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	d.enabled = temp.Enabled
	d.path = temp.Path
	return nil
}

func (d *Debug) MarshalJSON() ([]byte, error) {
	return json.Marshal(&debug{
		Enabled: d.enabled,
		Path:    d.path,
	})
}
