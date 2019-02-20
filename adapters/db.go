package adapters

import (
	"database/sql"
	"github.com/nikonm/go-gin-auth-middleware/user"
	"strings"
)

type DBProvider struct {
	db     *sql.DB
	config map[string]interface{}
}

/**
config = map[string]interface{}{
	"driver": "postgres",
	"connection": "postgres://user:pass@localhost/bookstore",
	"sql": "SELECT {select_columns} FROM users where username=$1 and password=$2",
	"source_target_fields": map[string]interface{}
}
*/
func (d *DBProvider) Init(config map[string]interface{}) error {
	d.config = config
	var err error

	d.db, err = d.connect()
	return err
}

func (d *DBProvider) SetCustomDb(db *sql.DB) {
	d.db = db
}

func (d *DBProvider) connect() (*sql.DB, error) {
	var err error
	if d.db != nil {
		if err = d.db.Ping(); err == nil {
			return d.db, nil
		}
		dbCopy := d.db
		go func(dc *sql.DB) {
			_ = dc.Close()
		}(dbCopy)
	}

	d.db, err = sql.Open(d.config["driver"].(string), d.config["connection"].(string))
	return d.db, err
}

func (d *DBProvider) Login(dto user.LoginDTO) (*user.User, error) {
	db, err := d.connect()
	if err != nil {
		return nil, err
	}

	cols := make([]string, 0)
	sqlSelect := ""

	for c := range d.config["source_target_fields"].(map[string]interface{}) {
		sqlSelect += c + ", "
		cols = append(cols, c)
	}
	sqlSelect = strings.TrimRight(sqlSelect, ", ")

	_sql := strings.Replace(d.config["sql"].(string), "{select_columns}", sqlSelect, 1)
	row := db.QueryRow(_sql, dto.User, dto.PasswordHash)

	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i, _ := range columns {
		columnPointers[i] = &columns[i]
	}

	// Scan the result into the column pointers...
	if err := row.Scan(columnPointers...); err != nil {
		return nil, err
	}

	// Create our map, and retrieve the value for each column from the pointers slice,
	// storing it in the map with the name of the column as the key.
	m := make(map[string]interface{})
	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		m[colName] = *val
	}

	return (&user.User{}).FillFromSource(m, d.config["source_target_fields"].(map[string]interface{}))
}
