package encryption_test

import (
	"goseed/utils/encryption"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparePasswords(t *testing.T) {
	password := []byte("MySecureL0ngPassw0rd")

	t.Run("success", func(t *testing.T) {
		hash := encryption.HashAndSalt(password)
		res := encryption.ComparePasswords(hash, password)
		assert.True(t, res)
	})

	t.Run("wrong-hash", func(t *testing.T) {
		hash := encryption.HashAndSalt([]byte("WrongPassword"))
		res := encryption.ComparePasswords(hash, password)
		assert.False(t, res)
	})
}
