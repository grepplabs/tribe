package masterkey

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestMasterKeySetCodec(t *testing.T) {
	a := assert.New(t)

	secret := []byte("testsecret")

	mks, err := NewMasterKeyset(secret)
	a.Nil(err)
	encryptedKey, err := mks.EncryptKeyset()
	a.Nil(err)

	mks2, err := DecryptKeyset(encryptedKey, secret)
	a.Nil(err)

	a.True(proto.Equal(mks.GetKeyset().KeysetInfo(), mks2.GetKeyset().KeysetInfo()), "key handlers are not equal")
}

func TestMasterKeySetWrongSecret(t *testing.T) {
	a := assert.New(t)

	mks, err := NewMasterKeyset([]byte("testsecret-1"))
	a.Nil(err)
	encryptedKey, err := mks.EncryptKeyset()
	a.Nil(err)

	_, err = DecryptKeyset(encryptedKey, []byte("testsecret-2"))
	a.NotNil(err, "decrypt with different password should fail")
}

func TestEmptyMasterSecret(t *testing.T) {
	a := assert.New(t)
	_, err := NewMasterKeyset([]byte{})
	a.Same(errEmptyMasterSecret, err)
}
