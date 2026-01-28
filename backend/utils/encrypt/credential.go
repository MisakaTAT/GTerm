package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"sync"

	"github.com/MisakaTAT/GTerm/backend/consts"
	"github.com/denisbrodbeck/machineid"
	"golang.org/x/crypto/pbkdf2"
)

const (
	iterationCount = 100000
	keyLength      = 32
	saltLength     = 16
)

var (
	once       sync.Once
	credential *Credential
	initErr    error
)

type Credential struct {
	machineID []byte
}

type Encrypted struct {
	Ciphertext string
	Salt       string
}

type CredentialField struct {
	Plaintext  string
	Ciphertext *string
	Salt       *string
}

func NewCredential() (*Credential, error) {
	once.Do(func() {
		id, err := machineid.ProtectedID(consts.ApplicationName)
		if err != nil {
			initErr = err
			return
		}
		credential = &Credential{
			machineID: []byte(id),
		}
	})
	return credential, initErr
}

func (c *Credential) EncryptPassword(plaintext string) (*Encrypted, error) {
	salt := make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key := pbkdf2.Key(c.machineID, salt, iterationCount, keyLength, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return &Encrypted{
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
		Salt:       base64.StdEncoding.EncodeToString(salt),
	}, nil
}

func (c *Credential) DecryptPassword(ciphertext, salt string) (string, error) {
	decodeSalt, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}

	key := pbkdf2.Key(c.machineID, decodeSalt, iterationCount, keyLength, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	decodeCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("insufficient ciphertext length")
	}
	nonce := decodeCiphertext[:gcm.NonceSize()]
	decodeCiphertext = decodeCiphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, decodeCiphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
