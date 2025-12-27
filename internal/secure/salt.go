package secure

import (
	"crypto/rand"
	"io"
	"wb_logistic_assistant/internal/errors"

	"wb_logistic_assistant/internal/utils"
)

func GetSalt(data []byte, from, to int) ([]byte, error) {
	res, err := utils.CutData(data, from, to)
	if err != nil {
		return nil, errors.Wrap(err, "Secure.GenerateSalt()", "")
	}
	return res, nil
}

func GetSaltReader(r io.Reader, from, to int) ([]byte, error) {
	res, err := utils.CutDataReader(r, from, to)
	if err != nil {
		return nil, errors.Wrap(err, "Secure.GetSaltReader()", "")
	}
	return res, nil
}

func GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, errors.Wrap(err, "Secure.GenerateSalt()", "failed to generate salt")
	}
	return salt, nil
}
