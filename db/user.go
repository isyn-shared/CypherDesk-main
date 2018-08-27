package db

import (
	"CypherDesk-main/alias"
	"database/sql"
	"fmt"
	"regexp"
)

var (
	findUserFields = []string{"Login", "Mail", "Name", "Surname", "Partonymic", "Recourse"}
)

// User - struct describing the registered user
type User struct {
	ID            int    `json: "id"`
	Login         string `json: "login"`
	Pass          string `json: "pass"`
	Mail          string `json: "mail"`
	Name          string `json: "name"`
	Surname       string `json: "surname"`
	Partonymic    string `json: "partonymic"`
	Recourse      string `json: "recourse"`
	Role          string `json: "role"`
	Department    int    `json: "department"`
	Status        string `json: "status"`
	ActivationKey string `json: "activationKey"`
}

type userNullFields struct {
	Mail          interface{}
	Name          interface{}
	Surname       interface{}
	Partonymic    interface{}
	Recourse      interface{}
	ActivationKey interface{}
	Status        interface{}
}

// BasicUser returns user obj containing system information
func BasicUser(mail, role, status string, department int) *User {
	user := &User{
		Mail:       mail,
		Role:       role,
		Status:     status,
		Department: department,
	}
	return user
}

// SetActivationKey
func (u *User) SetActivationKey(key string) {
	u.ActivationKey = alias.MD5(key)
}

// WriteIn fills empty fields fo user obj
func (u *User) WriteIn(user *User) {
	u.Name, u.Surname, u.Partonymic = user.Name, user.Surname, user.Partonymic
	u.Recourse = user.Recourse
	u.Login, u.Pass = user.Login, user.Pass

	if !alias.EmptyStr(user.ActivationKey) {
		u.ActivationKey = user.ActivationKey
	}
	if !alias.EmptyStr(user.Mail) {
		u.Mail = user.Mail
	}
}

// HashPass method encrypt password of user
func (u *User) HashPass() {
	u.Pass = alias.HashPass(u.Pass)
}

// UpdateUser fills empty fields of user entry in DB
func (m *MysqlUser) UpdateUser(user *User) int64 {
	db := m.connect()
	defer db.Close()
	stmt := prepare(db, "UPDATE users SET mail=?, name=?, surname=?, partonymic=?, recourse=?, login=?, pass=?, activationKey=? WHERE id = ?")
	defer stmt.Close()
	res := exec(stmt, []interface{}{user.Mail, user.Name, user.Surname, user.Partonymic,
		user.Recourse, user.Login, user.Pass, user.ActivationKey, user.ID})
	aff := affect(res)
	return aff
}

func (m *MysqlUser) InsertUser(user *User) sql.Result {
	db := m.connect()
	defer db.Close()
	stmt := prepare(db, "INSERT INTO users (login, pass, mail, name, surname, partonymic, recourse, role, department, status, activationKey) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	res := exec(stmt, []interface{}{user.Login, user.Pass, user.Mail, user.Name, user.Surname, user.Partonymic, user.Recourse,
		user.Recourse, user.Department, user.Status, user.ActivationKey})
	return res
}

// Exist method checks if user exist
func (u *User) Exist() bool {
	if u.ID == 0 {
		return false
	}
	return true
}

const loginCharset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateLogin generates login for user and check if user with this login exist`s
func (u *User) GenerateLogin(len int) {
	mysql := CreateMysqlUser()
	for {
		login := alias.StringWithCharset(len, loginCharset)
		if !mysql.GetUser("login", login).Exist() {
			u.Login = login
			return
		}
	}
}

// GeneratePass generates pass for user
func (u *User) GeneratePass(len int) {
	u.Pass = alias.StringWithCharset(len, loginCharset)
}

// Filled method returns true if all user fields are filled
func (u *User) Filled() bool {
	if u.Name == "" || u.Mail == "" {
		return false
	}
	return true
}

// GetUser return user obj using diff keys (id, login, mail)
func (m *MysqlUser) GetUser(sqlParam string, key interface{}) *User {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "SELECT * FROM users WHERE "+sqlParam+" = ? LIMIT 1")
	defer stmt.Close()

	user, ns := new(User), new(userNullFields)
	err := stmt.QueryRow(key).Scan(&user.ID, &user.Login, &user.Pass, &ns.Mail, &ns.Name, &ns.Surname,
		&ns.Partonymic, &ns.Recourse, &user.Role, &user.Department, &ns.Status, &ns.ActivationKey)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return user
	} else if err != nil {
		panic("db error: " + err.Error())
	}
	user.chkNullFields(ns)
	return user
}

// FindUser returns finded user-obj arr using key
func (m *MysqlUser) FindUser(keys []string) []*User {
	var argsLen = 6
	sqlReq := "SELECT * FROM users WHERE (name = ? OR surname = ? OR partonymic = ? OR login = ? OR mail = ? OR recourse = ?)"
	addSQLReq := " AND (name = ? OR surname = ? OR partonymic = ? OR login = ? OR mail = ? OR recourse = ?)"
	keysLen := len(keys)
	sqlKeys := make([]string, 0)

	for i := 0; i < argsLen; i++ {
		sqlKeys = append(sqlKeys, keys[0])
	}

	for i := 1; i < keysLen; i++ {
		sqlReq += addSQLReq
		for j := 0; j < argsLen; j++ {
			sqlKeys = append(sqlKeys, keys[i])
		}
	}

	fmt.Println(sqlReq)
	users := make([]*User, 0)

	db := m.connect()
	defer db.Close()

	stmt := prepare(db, sqlReq)
	defer stmt.Close()

	InterfaceArgs := make([]interface{}, len(sqlKeys))
	for i, a := range sqlKeys {
		InterfaceArgs[i] = a
	}

	rows := chk(stmt.Query(InterfaceArgs...)).(*sql.Rows)

	for rows.Next() {
		user, ns := new(User), new(userNullFields)
		err := rows.Scan(&user.ID, &user.Login, &user.Pass, &ns.Mail, &ns.Name, &ns.Surname,
			&ns.Partonymic, &ns.Recourse, &user.Role, &user.Department, &ns.Status, &ns.ActivationKey)
		if err != nil {
			fmt.Println(err.Error())
		}
		user.chkNullFields(ns)
		users = append(users, user)
	}

	return users
}

// TODO: think about refactoring
func (u *User) chkNullFields(ns *userNullFields) {
	if ns.Mail != nil {
		u.Mail = string(ns.Mail.([]byte))
	}
	if ns.Name != nil {
		u.Name = string(ns.Name.([]byte))
	}
	if ns.Surname != nil {
		u.Surname = string(ns.Surname.([]byte))
	}
	if ns.Partonymic != nil {
		u.Partonymic = string(ns.Partonymic.([]byte))
	}
	if ns.Recourse != nil {
		u.Recourse = string(ns.Recourse.([]byte))
	}
	if ns.ActivationKey != nil {
		u.ActivationKey = string(ns.ActivationKey.([]byte))
	}
	if ns.Status != nil {
		u.Status = string(ns.Status.([]byte))
	}
}

// GetDepartment method returns the department obj of user
func (u *User) GetDepartment() *Department {
	mysql := CreateMysqlUser()
	return mysql.GetDepartment(u.Department)
}

// HidePrivateInfo clears all private info from user obj
func (u *User) HidePrivateInfo() {
	u.Pass = ""
	u.ID = 0
}

// String returns user object description
func (u *User) String() string {
	return "Name: " + u.Name + "\nSurname: " + u.Surname + "\nPartonymic: " + u.Partonymic +
		"\nLogin: " + u.Login + "\nMail: " + u.Mail + "\nRecourse: " + u.Recourse
}

// ChkPass returns true if regexp match the password
func (u *User) ChkPass() bool {
	if alias.StrLen(u.Pass) > 5 && alias.StrLen(u.Pass) < 16 {
		match, _ := regexp.MatchString("^[a-zA-Z0-9]{1}[a-zA-Z0-9_-]+[a-zA-Z0-9]{1}$", u.Pass)
		return match
	}
	return false
}

// ChkLogin returns true if regexp match the login
func (u *User) ChkLogin() bool {
	if alias.StrLen(u.Login) > 3 && alias.StrLen(u.Login) < 11 {
		match, _ := regexp.MatchString("^[a-zA-Z0-9]{1}[a-zA-Z0-9_-]+[a-zA-Z0-9]{1}$", u.Login)
		return match
	}
	return false
}
