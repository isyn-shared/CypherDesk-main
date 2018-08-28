package db

import (
	"database/sql"
)

// Department obj
type Department struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetDepartment returns Department obj using id
func (m *MysqlUser) GetDepartment(sqlKey string, keyVal interface{}) *Department {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "SELECT * FROM departments WHERE "+sqlKey+" = ?")
	defer stmt.Close()

	d := new(Department)
	err := stmt.QueryRow(keyVal).Scan(&d.ID, &d.Name)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return d
	}
	if err != nil {
		panic("db error: " + err.Error())
	}
	return d
}

// InsertDepartment inserts
func (m *MysqlUser) InsertDepartment(name string) sql.Result {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "INSERT INTO departments (name) VALUES (?)")
	defer stmt.Close()
	res := exec(stmt, []interface{}{name})
	return res
}

func (m *MysqlUser) GetDepartmentUsers(sqlKey string, keyVal interface{}) {
	db := m.connect()
	defer db.Close()

	var depID int

	if sqlKey == "name" {
		dep := m.GetDepartment(sqlKey, keyVal)
		depID = dep.ID
	} else {
		depID = keyVal.(int)
	}

	stmt := prepare(db, "SELECT * FROM users WHERE department = ?")
	defer stmt.Close()

	rows := chk(stmt.Query(depID)).(*sql.Rows)

	for rows.Next() {
		//TODO: Дописать!!!!!!!!!!!
	}
}

// GetDepartments return all departmens objects from DB
func (m *MysqlUser) GetDepartments() []*Department {
	db := m.connect()
	defer db.Close()

	rows := getQuery(db, "SELECT * FROM departments")
	res := make([]*Department, 0)
	for rows.Next() {
		dep := new(Department)
		err := rows.Scan(&dep.ID, &dep.Name)
		if err != nil {
			panic("db error: " + err.Error())
		}
		res = append(res, dep)
	}
	return res
}
