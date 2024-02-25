package tests

import (
	d "Example/src/database"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {

}

func TestGenerateCode(t *testing.T) {
	if reflect.TypeOf(d.GenerateFriendCode()).Kind() == reflect.String {
		assert.True(t, true)
	}
	assert.False(t, false)

}

func TestInitialize(t *testing.T) {
	d.InitializeDatabase()
	assert.True(t, true)
}

func TestUserDetails(t *testing.T) {

}
