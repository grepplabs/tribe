package masterkey

import (
	"bytes"
	"crypto/sha256"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/aead/subtle"
	"github.com/google/tink/go/keyset"
	"github.com/pkg/errors"
)

var errEmptyMasterSecret = errors.New("kms: master secret is empty")

type MasterKeyset interface {
	EncryptKeyset() ([]byte, error)
	GetKeyset() *keyset.Handle
}

type masterKeyset struct {
	kh     *keyset.Handle
	secret []byte
}

func NewMasterKeyset(secret []byte) (MasterKeyset, error) {
	if len(secret) == 0 {
		return nil, errEmptyMasterSecret
	}
	kh, err := keyset.NewHandle(aead.AES256CTRHMACSHA256KeyTemplate())
	if err != nil {
		return nil, err
	}
	return &masterKeyset{
		secret: secret,
		kh:     kh,
	}, nil
}

// DecryptKeyset decrypts encrypted key set using provided secret
func DecryptKeyset(encryptedKeyset []byte, secret []byte) (MasterKeyset, error) {
	masterKey, err := getKMSEnvelopeAEAD(secret)
	if err != nil {
		return nil, err
	}
	r := keyset.NewBinaryReader(bytes.NewBuffer(encryptedKeyset))
	kh, err := keyset.Read(r, masterKey)
	if err != nil {
		return nil, err
	}
	return &masterKeyset{
		kh:     kh,
		secret: secret,
	}, nil
}

// EncryptKeyset encrypts key set
func (m masterKeyset) EncryptKeyset() ([]byte, error) {
	masterKey, err := getKMSEnvelopeAEAD(m.secret)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	w := keyset.NewBinaryWriter(buf)
	if err := m.kh.Write(w, masterKey); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m masterKeyset) GetKeyset() *keyset.Handle {
	return m.kh
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
