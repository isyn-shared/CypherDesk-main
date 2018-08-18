package db

import (
	"database/sql"
	"log"

	// Wrapper for sql package
	_ "github.com/go-sql-driver/mysql"
)

// NullString - obj from sql package
type NullString sql.NullString

func chk(obj interface{}, err error) interface{} {
	if err != nil {
		log.Fatal("panic in db: " + err.Error())
		panic(err.Error())
	}
	return obj
}

func getQuery(db *sql.DB, sqlReq string) *sql.Rows {
	return chk(db.Query(sqlReq)).(*sql.Rows)
}

func prepare(db *sql.DB, sqlReq string) *sql.Stmt {
	return chk(db.Prepare(sqlReq)).(*sql.Stmt)
}

func exec(stmt *sql.Stmt, args []interface{}) sql.Result {
	return chk(stmt.Exec(args...)).(sql.Result)
}

func affect(res sql.Result) int64 {
	return chk(res.RowsAffected()).(int64)
}
