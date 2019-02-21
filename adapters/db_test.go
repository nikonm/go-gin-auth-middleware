package adapters

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nikonm/go-gin-auth-middleware/user"
	"testing"
)

func getDBMock(t *testing.T) *sql.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	rows := sqlmock.NewRows([]string{
		"id", "username", "email"}).
		AddRow("1", "test", "test@test.local")
	mock.ExpectQuery("SELECT (.*) FROM users where username=\\? and password=\\?").WithArgs("test", "test").WillReturnRows(rows)

	return db
}

func getDbProvider(t *testing.T) DBProvider {
	p := DBProvider{}
	db := getDBMock(t)
	p.SetCustomDb(db)
	config := map[string]interface{}{
		"driver":     "postgres",
		"connection": "postgres://user:pass@localhost/test",
		"sql":        "SELECT {select_columns} FROM users where username=? and password=?",
		"source_target_fields": map[string]interface{}{
			"id":       "id",
			"username": "user",
			"email":    "email",
		},
	}
	err := p.Init(config)
	if err != nil {
		t.Fail()
	}
	return p
}

func TestDBProvider_Login(t *testing.T) {
	p := getDbProvider(t)
	dto := user.LoginDTO{User: "test", Password: "test", Type: "db"}
	dto.PasswordHash = "test"
	u, err := p.Login(dto)
	//t.Log(u, err)
	if u == nil {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
}
