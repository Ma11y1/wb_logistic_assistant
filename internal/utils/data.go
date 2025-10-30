package utils

import (
	"errors"
	"fmt"
	"io"
)

func CutData(data []byte, from, to int) ([]byte, error) {
	if from < 0 || to < 0 || from > to {
		return nil, fmt.Errorf("invalid length range: from %d to %d", from, to)
	}
	if from == to {
		return []byte{}, nil
	}
	if len(data) < to {
		return nil, fmt.Errorf("data length must be %d bytes more", to)
	}
	return data[from:to], nil
}

func CutDataReader(reader io.Reader, from, to int) ([]byte, error) {
	if from < 0 || to < 0 || from > to {
		return nil, fmt.Errorf("invalid range: from %d to %d", from, to)
	}
	if from == to {
		return []byte{}, nil
	}
	if reader == nil {
		return nil, fmt.Errorf("nil reader")
	}

	if seeker, ok := reader.(io.Seeker); ok {
		_, err := seeker.Seek(int64(from), io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("seek error: %w", err)
		}
	} else {
		_, err := io.CopyN(io.Discard, reader, int64(from))
		if err != nil {
			return nil, fmt.Errorf("discard error: %w", err)
		}
	}

	data := make([]byte, to-from)
	n, err := io.ReadFull(reader, data)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, fmt.Errorf("read error: %w", err)
	}
	return data[:n], nil
}
