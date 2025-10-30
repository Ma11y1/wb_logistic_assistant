package secure

import (
	"crypto/rand"
	"sync"
)

type SecurePassword struct {
	parts  [4][]byte
	key    []byte
	mtx    sync.RWMutex
	length int
}

func NewSecurePasswordPure() (*SecurePassword, error) {
	sp := &SecurePassword{}

	sp.key = make([]byte, 16)
	if _, err := rand.Read(sp.key); err != nil {
		return nil, err
	}

	for i := range sp.parts {
		sp.parts[i] = make([]byte, 0, 32)
	}

	return sp, nil
}

func NewSecurePassword(password []byte) (*SecurePassword, error) {
	sp := &SecurePassword{length: len(password)}

	sp.key = make([]byte, 16)
	if _, err := rand.Read(sp.key); err != nil {
		return nil, err
	}

	for i := range sp.parts {
		sp.parts[i] = make([]byte, 0, 32)
	}

	sp.setPassword(password)
	return sp, nil
}

func (p *SecurePassword) Set(password []byte) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.wipe()
	p.length = len(password)
	p.setPassword(password)
}

func (p *SecurePassword) Get() string {
	p.mtx.RLock()
	defer p.mtx.RUnlock()
	return string(p.getPassword())
}

func (p *SecurePassword) Bytes() []byte {
	p.mtx.RLock()
	defer p.mtx.RUnlock()
	return p.getPassword()
}

func (p *SecurePassword) Wipe() {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.wipe()
}

func (p *SecurePassword) Len() int {
	p.mtx.RLock()
	defer p.mtx.RUnlock()
	return p.length
}

func (p *SecurePassword) setPassword(password []byte) {
	if len(password) == 0 {
		return
	}

	pwdBytes := password

	encrypted := make([]byte, len(pwdBytes))
	for i := range pwdBytes {
		encrypted[i] = pwdBytes[i] ^ p.key[i%len(p.key)]
	}

	partSize := (len(encrypted) + 3) / 4
	for i := 0; i < 4; i++ {
		start := i * partSize
		end := start + partSize
		if end > len(encrypted) {
			end = len(encrypted)
		}
		if start < len(encrypted) {
			noise := make([]byte, 8)
			_, _ = rand.Read(noise)
			p.parts[i] = append(p.parts[i], noise...)
			p.parts[i] = append(p.parts[i], encrypted[start:end]...)
		}
	}
}

func (p *SecurePassword) getPassword() []byte {
	if p.length == 0 {
		return []byte{}
	}

	var encrypted []byte
	for i := 0; i < 4; i++ {
		if len(p.parts[i]) > 8 { // Skip the first 8 bytes of noise
			encrypted = append(encrypted, p.parts[i][8:]...)
		}
	}

	if len(encrypted) > p.length {
		encrypted = encrypted[:p.length]
	}

	decrypted := make([]byte, len(encrypted))
	for i := range encrypted {
		decrypted[i] = encrypted[i] ^ p.key[i%len(p.key)]
	}

	return decrypted
}

func (p *SecurePassword) wipe() {
	for i := range p.parts {
		if len(p.parts[i]) > 0 {
			_, _ = rand.Read(p.parts[i])
			for j := range p.parts[i] {
				p.parts[i][j] = 0
			}
			p.parts[i] = p.parts[i][:0]
		}
	}
	p.length = 0
}
