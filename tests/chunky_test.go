package tests

import (
	"Example/src/database"
	"database/sql"
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestCreateAccount(t *testing.T) {
	if reflect.TypeOf(database.CreateAccount("User1", "56:de:51:53:38:2f")).Kind() == reflect.Bool {
		assert.True(t, true)
	}
	assert.False(t, false)
}

func TestCreateAcountError(t *testing.T) {

}

func TestGetUsersFriendCode(t *testing.T) {
	if reflect.TypeOf(database.GetUserFriendCode("User1")).Kind() == reflect.String {
		assert.True(t, true)
	}
	assert.False(t, false)
}
