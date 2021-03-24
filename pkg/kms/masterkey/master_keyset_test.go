package masterkey

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestMasterKeySetCodec(t *testing.T) {
	a := assert.New(t)

	secret := []byte("testsecret")

	mks, err := NewMasterKeySet(secret)
	a.Nil(err)
	encryptedKey, err := mks.EncryptKeySet()
	a.Nil(err)

	mks2, err := Decrypt(encryptedKey, secret)
	a.Nil(err)

	a.True(proto.Equal(mks.GetHandle().KeysetInfo(), mks2.GetHandle().KeysetInfo()), "key handlers are not equal")
}

func TestMasterKeySetWrongSecret(t *testing.T) {
	a := assert.New(t)

	mks, err := NewMasterKeySet([]byte("testsecret-1"))
	a.Nil(err)
	encryptedKey, err := mks.EncryptKeySet()
	a.Nil(err)

	_, err = Decrypt(encryptedKey, []byte("testsecret-2"))
	a.NotNil(err, "decrypt with different password should fail")
}

func TestEmptyMasterSecret(t *testing.T) {
	a := assert.New(t)
	_, err := NewMasterKeySet([]byte{})
	a.Same(errEmptyMasterSecret, err)
}
