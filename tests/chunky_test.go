package tests

/*
Uncomment test code at your own peril
*/

import (
	"database/sql"
	"log"

	sm "github.com/DATA-DOG/go-sqlmock"
)

func NewMock() (*sql.DB, sm.Sqlmock) {
	db, mock, err := sm.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

// func TestCreateAccount(t *testing.T) {
// 	//setup
// 	var errExp error
// 	errExp = nil

// 	//invoke
// 	errAct := database.CreateAccount("user1", "56:de:51:53:38:2f")

// 	//analyze
// 	if errExp == errAct {
// 		assert.True(t, true)
// 	}
// 	assert.False(t, false)
// }

// func TestGetUsersFriendCode(t *testing.T) {
// 	//setup
// 	var errExp error
// 	errExp = nil

// 	//invoke
// 	_, errAct := database.GetUserFriendCode("User")

// 	//analyze
// 	if errExp == errAct {
// 		assert.True(t, true)
// 	}
// 	assert.False(t, false)
// }
