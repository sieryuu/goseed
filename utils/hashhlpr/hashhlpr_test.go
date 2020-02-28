package hashhlpr_test

import (
	"goseed/utils/hashhlpr"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparePasswords(t *testing.T) {
	password := []byte("MySecureL0ngPassw0rd")

	t.Run("success", func(t *testing.T) {
		hashString := hashhlpr.HashAndSalt(password)
		res := hashhlpr.ComparePasswords(hashString, password)
		assert.True(t, res)
	})

	t.Run("wrong-hash", func(t *testing.T) {
		hashString := hashhlpr.HashAndSalt([]byte("WrongPassword"))
		res := hashhlpr.ComparePasswords(hashString, password)
		assert.False(t, res)
	})
}
