package storage

import (
	"bytes"
	"io"
	"os"
	"sync"
	"wb_logistic_assistant/internal/config"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/secure"
)

type FileStorage struct {
	mtx          sync.Mutex
	config       *config.Config
	storageModel *storageModel
	serializer   Serializer
	crypter      secure.Crypter
	key          secure.Key
	password     *secure.SecurePassword
	isEncrypted  bool
}

func NewFileStorage(config *config.Config) (*FileStorage, error) {
	password, err := secure.NewSecurePasswordPure()
	if err != nil {
		return nil, errors.New("FileStorage.New()", "failed to initialize secure password")
	}
	return &FileStorage{
		config: config,
		storageModel: &storageModel{
			ConfigStore: NewFileConfigStore(),
			CacheStore:  NewFileCache(),
		},
		serializer:  JSONSerializer{},
		password:    password,
		isEncrypted: false,
	}, nil
}

func (s *FileStorage) ConfigStore() ConfigStore {
	return s.storageModel.ConfigStore
}

func (s *FileStorage) CacheStore() CacheStore {
	return s.storageModel.CacheStore
}

func (s *FileStorage) Load(path string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrapf(err, "FileStorage.Load()", "failed to open file storage by path %s", path)
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return errors.Wrapf(err, "FileStorage.Load()", "failed to read file storage by path %s", path)
	}

	// If the file is missing or empty, then we don’t load anything.
	if len(fileData) < 1 {
		return nil
	}

	encryptedMarker := fileData[0]
	fileData = fileData[1:]

	if s.isEncrypted && encryptedMarker == encryptStorageMarker {
		fileData, err = s.crypter.Decrypt(fileData)
		if err != nil {
			return errors.Wrapf(err, "FileStorage.Load()", "failed to decrypt file storage by path %s", path)
		}
	} else if !s.isEncrypted && encryptedMarker == encryptStorageMarker {
		return errors.Newf("FileStorage.Load()", "file %s is encrypted, but the application has not been given a decryption key", path)
	}

	if err = s.serializer.Decode(bytes.NewBuffer(fileData), s.storageModel); err != nil && !errors.Is(err, io.EOF) {
		return errors.New("FileStorage.Load()", "failed to decode storage model")
	}

	return nil
}

func (s *FileStorage) Save(path string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	var buf bytes.Buffer
	err := s.serializer.Encode(&buf, s.storageModel)
	if err != nil {
		return errors.Wrap(err, "FileStorage.Save()", "failed to encode storage model")
	}

	data := buf.Bytes()
	marker := rawStorageMarker

	if s.isEncrypted {
		marker = encryptStorageMarker
		data, err = s.crypter.Encrypt(data)
		if err != nil {
			return errors.Wrap(err, "FileStorage.Save()", "failed to encrypt storage model")
		}
	}

	tmpPath := path + ".tmp"
	file, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrapf(err, "FileStorage.Save()", "failed to open temporary file %s", tmpPath)
	}
	defer file.Close()

	if _, err = file.Write(append([]byte{marker}, data...)); err != nil {
		return errors.Wrapf(err, "FileStorage.Save()", "failed to write temporary file %s", tmpPath)
	}
	if err = file.Sync(); err != nil {
		return errors.Wrapf(err, "FileStorage.Save()", "failed to sync temporary file\n %s", tmpPath)
	}
	if err = file.Close(); err != nil {
		return errors.Wrapf(err, "FileStorage.Save()", "failed to close temporary file\n %s", tmpPath)
	}

	if err = os.Rename(tmpPath, path); err != nil {
		return errors.Wrapf(err, "FileStorage.Save()", "failed to rename temporary file\n %s в %s", tmpPath, path)
	}

	return nil
}

func (s *FileStorage) Clear() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.storageModel.ConfigStore.Clear()
	s.storageModel.CacheStore.Clear()
}

// SetEncrypt Sets the encryption password and enables storage encryption, if the password is empty encryption is disabled
func (s *FileStorage) SetEncrypt(password []byte) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.password.Wipe()
	s.isEncrypted = len(password) > 0

	if !s.isEncrypted {
		s.key = nil
		s.crypter = nil
		return
	}

	s.password.Set(password)
	s.key = secure.NewArgon2IDKey(
		s.password,
		s.config.Storage().ArgonTime(),
		s.config.Storage().ArgonMemory(),
		s.config.Storage().ArgonKeyLength(),
		s.config.Storage().ArgonThreads(),
	)

	s.crypter = &secure.AESGCMCrypter{
		SaltSize: s.config.Storage().SaltSize(),
		Key:      s.key,
	}

	for i := 0; i < len(password); i++ {
		password[i] = '0'
	}
}

func (s *FileStorage) IsEncrypted() bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.isEncrypted
}
