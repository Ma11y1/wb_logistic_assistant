package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"wb_logistic_assistant/internal/errors"
)

type Encrypter interface {
	Encrypt([]byte) ([]byte, error)
}

type Decrypter interface {
	Decrypt([]byte) ([]byte, error)
}

type Crypter interface {
	Encrypter
	Decrypter
}

type AESGCMCrypter struct {
	SaltSize int
	Key      Key
}

// Encrypt Salt property of the Key structure can be modified
func (c *AESGCMCrypter) Encrypt(data []byte) ([]byte, error) {
	if c.SaltSize > 0 {
		salt, err := GenerateSalt(c.SaltSize)
		if err != nil {
			return nil, errors.Wrap(err, "AESGCMCrypter.encrypt()", "failed to generate salt")
		}
		c.Key.SetSalt(salt)
	}

	key, err := c.Key.Get()
	if err != nil {
		return nil, errors.Wrap(err, "AESGCMCrypter.encrypt()", "failed to generate key")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "AESGCMCrypter.encrypt()", "failed to create cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "AESGCMCrypter.encrypt()", "failed to create GCM")
	}

	iv := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(iv); err != nil {
		return nil, errors.Wrap(err, "AESGCMCrypter.encrypt()", "failed to generate IV")
	}

	ciphertext := gcm.Seal(nil, iv, data, nil)

	result := make([]byte, 0, len(ciphertext))
	if salt := c.Key.GetSalt(); salt != nil {
		result = append(result, salt...)
	}
	result = append(result, iv...)
	result = append(result, ciphertext...)

	return result, nil
}

// Decrypt Salt property of the Key structure can be modified
func (c *AESGCMCrypter) Decrypt(data []byte) ([]byte, error) {
	ivSize := 12
	minLen := c.SaltSize + ivSize + 1
	if len(data) < minLen {
		return nil, errors.New("AESGCMCrypter.Decrypt()", "input too short")
	}

	if c.SaltSize > 0 {
		salt := data[:c.SaltSize]
		c.Key.SetSalt(salt)
	}

	iv := data[c.SaltSize : c.SaltSize+ivSize]
	ciphertext := data[c.SaltSize+ivSize:]

	key, err := c.Key.Get()
	if err != nil {
		return nil, errors.Wrap(err, "AESGCMCrypter.Decrypt()", "failed to derive key after salt")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "AESGCMCrypter.Decrypt()", "failed to create cipher block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "AESGCMCrypter.Decrypt()", "failed to create GCM")
	}

	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, errors.Wrap(err, "AESGCMCrypter.Decrypt()", "decryption failed")
	}

	return plaintext, nil
}
