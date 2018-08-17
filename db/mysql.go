package db

import (
	"CypherDesk-main/alias"
	"database/sql"
	"strings"

	// Wrapper for sql package
	_ "github.com/go-sql-driver/mysql"
)

// MysqlUser - obj that can connect to mysql database and make requests
type MysqlUser struct {
	user     string
	password string
	database string
}

func (m *MysqlUser) create(user, pass, db string) {
	m.user = user
	m.password = pass
	m.database = db
}

func (m MysqlUser) connect() *sql.DB {
	var login string
	login = m.user + ":" + m.password + "@/" + m.database + "?parseTime=true"
	return chk(sql.Open("mysql", login)).(*sql.DB)
}

//CreateMysqlUser returns new authorized mysql user
func CreateMysqlUser() *MysqlUser {
	mysql := new(MysqlUser)
	sqlLogin, sqlPass, dbName := getMysqlKey()
	mysql.create(sqlLogin, sqlPass, dbName)
	return mysql
}

func getMysqlKey() (string, string, string) {
	bs := chk(alias.ReadFile("mysql.key"))

	str := bs.(string)
	lp := strings.Split(str, ";")

	if len(lp) != 3 {
		panic("Неправильный формат файла mysql.key!")
	}

	return lp[0], lp[1], lp[2]
}
