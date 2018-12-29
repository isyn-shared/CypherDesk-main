package db

import (
	"database/sql"
	"fmt"
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

// GetDepartmentUsers returns all users which con
func (m *MysqlUser) GetDepartmentUsers(sqlKey string, keyVal interface{}) []*User {
	db := m.connect()
	defer db.Close()

	var depID int

	if sqlKey == "name" {
		dep := m.GetDepartment(sqlKey, keyVal)
		depID = dep.ID
	} else {
		depID = keyVal.(int)
	}

	stmt := prepare(db, "SELECT * FROM users WHERE department = ? AND role = 'mlRzUQ=='")
	defer stmt.Close()

	rows := chk(stmt.Query(depID)).(*sql.Rows)

	users := make([]*User, 0)
	for rows.Next() {
		user, ns := new(User), new(userNullFields)
		err := rows.Scan(&user.ID, &user.Login, &user.Pass, &ns.Mail, &ns.Name, &ns.Surname,
			&ns.Partonymic, &ns.Recourse, &user.Role, &user.Department, &ns.Status, &ns.ActivationKey, &ns.ActivationType)
		if err != nil {
			fmt.Println(err.Error())
		}
		user.chkNullFields(ns)
		user.RefactStandartInfo(true)
		users = append(users, user)
	}

	return users
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

// UpdateDepartment updates department entry in the db
func (m *MysqlUser) UpdateDepartment(oldDep *Department, newName string) int64 {
	db := m.connect()
	defer db.Close()

	var sqlKey string
	var keyVal interface{}
	if oldDep.ID == 0 {
		sqlKey = "name"
		keyVal = oldDep.Name
	} else {
		sqlKey = "id"
		keyVal = oldDep.ID
	}

	stmt := prepare(db, "UPDATE departments SET name=? WHERE "+sqlKey+"=?")
	defer stmt.Close()

	res := exec(stmt, []interface{}{newName, keyVal})
	aff := affect(res)
	return aff
}
