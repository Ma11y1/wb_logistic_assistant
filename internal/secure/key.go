package secure

import (
	"golang.org/x/crypto/argon2"
	"wb_logistic_assistant/internal/errors"
)

type Key interface {
	Get() ([]byte, error)
	SetSalt(salt []byte)
	GetSalt() []byte
}

type Argon2IDKey struct {
	password *SecurePassword
	salt     []byte
	time     uint32
	memory   uint32
	keyLen   uint32
	threads  uint8
}

func NewArgon2IDKey(password *SecurePassword, time, memory, keyLen uint32, threads uint8) *Argon2IDKey {
	return &Argon2IDKey{
		password: password,
		time:     time,
		memory:   memory,
		keyLen:   keyLen,
		threads:  threads,
	}
}

func (k *Argon2IDKey) Get() ([]byte, error) {
	if k.password == nil || k.password.Len() == 0 {
		return nil, errors.Wrap(nil, "Argon2IDKey.Get()", "password is nil or empty")
	}
	return argon2.IDKey(k.password.Bytes(), k.salt, k.time, k.memory, k.threads, k.keyLen), nil
}

func (k *Argon2IDKey) SetSalt(salt []byte) {
	k.salt = salt
}

func (k *Argon2IDKey) GetSalt() []byte {
	return k.salt
}
