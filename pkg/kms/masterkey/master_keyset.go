package masterkey

import (
	"bytes"
	"crypto/sha256"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/aead/subtle"
	"github.com/google/tink/go/keyset"
	"github.com/pkg/errors"
	"log"
)

var errEmptyMasterSecret = errors.New("kms: master secret is empty")

type MasterKeySet interface {
	EncryptKeySet() ([]byte, error)
	GetHandle() *keyset.Handle
	KeyId() uint32
}

type masterKeySet struct {
	kh     *keyset.Handle
	secret []byte
}

func NewMasterKeySet(secret []byte) (MasterKeySet, error) {
	if len(secret) == 0 {
		return nil, errEmptyMasterSecret
	}
	kh, err := keyset.NewHandle(aead.AES256CTRHMACSHA256KeyTemplate())
	if err != nil {
		return nil, err
	}
	return &masterKeySet{
		secret: secret,
		kh:     kh,
	}, nil
}

// Decrypt decrypts encrypted key set using provided secret
func Decrypt(encryptedKeySet []byte, secret []byte) (MasterKeySet, error) {
	masterKey, err := getKMSEnvelopeAEAD(secret)
	if err != nil {
		return nil, err
	}
	r := keyset.NewBinaryReader(bytes.NewBuffer(encryptedKeySet))
	kh, err := keyset.Read(r, masterKey)
	if err != nil {
		return nil, err
	}
	return &masterKeySet{
		kh:     kh,
		secret: secret,
	}, nil
}

// EncryptKeySet encrypts key set
func (m masterKeySet) EncryptKeySet() ([]byte, error) {
	masterKey, err := getKMSEnvelopeAEAD(m.secret)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	w := keyset.NewBinaryWriter(buf)
	if err := m.kh.Write(w, masterKey); err != nil {
		log.Fatal(err)
	}
	return buf.Bytes(), nil
}

func (m masterKeySet) GetHandle() *keyset.Handle {
	return m.kh
}

func (m masterKeySet) KeyId() uint32 {
	return m.kh.KeysetInfo().PrimaryKeyId
}

func getKMSEnvelopeAEAD(secret []byte) (*aead.KMSEnvelopeAEAD, error) {
	key := hashByteSecret(secret)
	backend, err := subtle.NewAESGCM(key)
	if err != nil {
		return nil, err
	}
	return aead.NewKMSEnvelopeAEAD2(aead.AES256GCMKeyTemplate(), backend), nil
}

// 32 bytes for AES-256
func hashByteSecret(secret []byte) []byte {
	r := sha256.Sum256(secret)
	return r[:]
}
