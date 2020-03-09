package hashutil_test

import (
	"goseed/utils/hashutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparePasswords(t *testing.T) {
	password := []byte("MySecureL0ngPassw0rd")

	t.Run("success", func(t *testing.T) {
		hashString := hashutil.HashAndSalt(password)
		res := hashutil.ComparePasswords(hashString, password)
		assert.True(t, res)
	})

	t.Run("wrong-hash", func(t *testing.T) {
		hashString := hashutil.HashAndSalt([]byte("WrongPassword"))
		res := hashutil.ComparePasswords(hashString, password)
		assert.False(t, res)
	})
}
