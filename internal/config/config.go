package config

import (
	"encoding/json"
	"os"
	"wb_logistic_assistant/internal/errors"
)

type Config struct {
	debug        *Debug        // ro
	internal     *Internal     // ro
	reports      *Reports      // ro
	storage      *Storage      // ro
	logistic     *Logistic     // ro
	googleSheets *GoogleSheets // ro
	telegram     *TelegramBot  // ro
}

type config struct {
	Debug        *Debug        `json:"debug"`
	Reports      *Reports      `json:"reports"`
	Storage      *Storage      `json:"storage"`
	Logistic     *Logistic     `json:"logistic"`
	GoogleSheets *GoogleSheets `json:"google_sheets"`
	Telegram     *TelegramBot  `json:"telegram_bot"`
}

func NewConfigFile(filePath string) (*Config, error) {
	c := &Config{
		debug:        newDebug(),        // default
		internal:     newInternal(),     // default
		reports:      newReports(),      // default
		storage:      newStorage(),      // default
		logistic:     newLogistic(),     // default
		googleSheets: newGoogleSheets(), // default
		telegram:     newTelegramBot(),  // default
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "Config.New()", "Failed opening config file by path %s", filePath)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(c)
	if err != nil {
		return nil, errors.Wrapf(err, "Config.New()", "Failed decoding config file by path %s", filePath)
	}

	err = validation(c)
	if err != nil {
		return nil, errors.Wrap(err, "Config.New()", "Config is invalid")
	}

	return c, nil
}

func (c *Config) Debug() *Debug               { return c.debug }
func (c *Config) Internal() *Internal         { return c.internal }
func (c *Config) Reports() *Reports           { return c.reports }
func (c *Config) Storage() *Storage           { return c.storage }
func (c *Config) Logistic() *Logistic         { return c.logistic }
func (c *Config) GoogleSheets() *GoogleSheets { return c.googleSheets }
func (c *Config) Telegram() *TelegramBot      { return c.telegram }

func (c *Config) UnmarshalJSON(b []byte) error {
	temp := &config{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	c.debug = temp.Debug
	c.reports = temp.Reports
	c.storage = temp.Storage
	c.googleSheets = temp.GoogleSheets
	c.logistic = temp.Logistic
	c.telegram = temp.Telegram
	return nil
}

func (c *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(&config{
		Debug:        c.debug,
		Reports:      c.reports,
		Storage:      c.storage,
		Logistic:     c.logistic,
		GoogleSheets: c.googleSheets,
		Telegram:     c.telegram,
	})
}
