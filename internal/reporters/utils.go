package reporters

import (
	"context"
	"strconv"
	"time"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/logger"
)

var NullTime = time.Time{}

func atoiSafe(s string) int {
	v, _ := strconv.ParseInt(s, 10, 64)
	return int(v)
}

func atofSafe(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return f
}

func retryAction(ctx context.Context, source string, attempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		logger.Logf(logger.WARN, "Reporters.retryAction()", "failed action %s, attempt %d/%d: %v", source, i+1, attempts, err)

		delay *= time.Duration(1 << i)
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "Reporters.retryAction()", "context cancelled while retrying")
		}
	}
	return errors.Wrapf(err, "Reporters.retryAction()", "all %d attempts failed for %s", attempts, source)
}

// SpNameToGateParking parse sp name and returns values.
//
//	Input example: шуш143.вор89-92.буфер303, Буфер_Маршрут 62, СЦМШ143_Буфер в МШ 177
//	Returns example: 89, 303
func SpNameToGateParking(s string) (gate, parking int) {
	b := []byte(s)
	n := len(b)

	i := 0
	for i < n {

		// "вор"
		if i+6 <= n &&
			b[i] == 0xD0 && b[i+1] == 0xB2 &&
			b[i+2] == 0xD0 && b[i+3] == 0xBE &&
			b[i+4] == 0xD1 && b[i+5] == 0x80 {

			i += 6
			gate = 0
			for i < n && b[i] == ' ' {
				i++
			}

			for i < n && b[i] >= '0' && b[i] <= '9' {
				gate = gate*10 + int(b[i]-'0')
				i++
			}
			continue
		}

		// "буфер"
		if i+10 <= n &&
			b[i] == 0xD0 && b[i+1] == 0xB1 &&
			b[i+2] == 0xD1 && b[i+3] == 0x83 &&
			b[i+4] == 0xD1 && b[i+5] == 0x84 &&
			b[i+6] == 0xD0 && b[i+7] == 0xB5 &&
			b[i+8] == 0xD1 && b[i+9] == 0x80 {

			i += 10
			parking = 0
			for i < n && b[i] == ' ' {
				i++
			}

			for i < n && b[i] >= '0' && b[i] <= '9' {
				parking = parking*10 + int(b[i]-'0')
				i++
			}
			continue
		}

		if i+5 <= n && b[i] == ' ' &&
			b[i+1] == 0xD0 && b[i+2] == 0x9C && // М
			b[i+3] == 0xD0 && b[i+4] == 0xA8 { // Ш

			i += 5
			parking = 0
			for i < n && b[i] == ' ' {
				i++
			}
			for i < n && b[i] >= '0' && b[i] <= '9' {
				parking = parking*10 + int(b[i]-'0')
				i++
			}
			continue
		}

		// "Буфер_Маршрут", "Маршрут"
		if i+14 <= n &&
			b[i] == 0xD0 && b[i+1] == 0x9C && // М
			b[i+2] == 0xD0 && b[i+3] == 0xB0 && // а
			b[i+4] == 0xD1 && b[i+5] == 0x80 && // р
			b[i+6] == 0xD1 && b[i+7] == 0x88 && // ш
			b[i+8] == 0xD1 && b[i+9] == 0x80 && // р
			b[i+10] == 0xD1 && b[i+11] == 0x83 && // у
			b[i+12] == 0xD1 && b[i+13] == 0x82 { // т

			i += 14
			parking = 0
			for i < n && b[i] == ' ' {
				i++
			}

			for i < n && b[i] >= '0' && b[i] <= '9' {
				parking = parking*10 + int(b[i]-'0')
				i++
			}
			continue
		}

		i++
	}

	return
}
