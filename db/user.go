package db

import (
	"CypherDesk-main/alias"
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
)

var (
	findUserFields = []string{"Login", "Mail", "Name", "Surname", "Partonymic", "Recourse"}
)

const (
	stInfoKey = "keys/userdatakey.toml"
	passKey   = "keys/passkey.toml"
)

// User - struct describing the registered user
type User struct {
	ID             int    `json: "id"`
	Login          string `json: "login"`
	Pass           string `json: "pass"`
	Mail           string `json: "mail"`
	Name           string `json: "name"`
	Surname        string `json: "surname"`
	Partonymic     string `json: "partonymic"`
	Recourse       string `json: "recourse"`
	Role           string `json: "role"`
	Department     int    `json: "department"`
	Status         string `json: "status"`
	ActivationKey  string `json: "activationKey"`
	ActivationType int    `json: "activationType"`
}

type userNullFields struct {
	Mail           interface{}
	Name           interface{}
	Surname        interface{}
	Partonymic     interface{}
	Recourse       interface{}
	ActivationKey  interface{}
	Status         interface{}
	ActivationType interface{}
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
	u.ActivationType = 1
}

func (u *User) SetRemindKey(key string) {
	u.ActivationKey = alias.MD5(key)
	u.ActivationType = 2
}

// WriteIn fills empty fields fo user obj
func (u *User) WriteIn(user *User) {
	u.Name, u.Surname, u.Partonymic = user.Name, user.Surname, user.Partonymic
	u.Recourse = user.Recourse
	u.Login, u.Pass = user.Login, user.Pass

	if !alias.EmptyStr(user.ActivationKey) {
		u.ActivationKey = user.ActivationKey
	}
	if user.ActivationType == 0 {
		u.ActivationType = user.ActivationType
	}
	if !alias.EmptyStr(user.Mail) {
		u.Mail = user.Mail
	}
}

// RefactStandartInfo encrypts/decrypts all string-fields of user
func (u *User) RefactStandartInfo(dec bool) {
	ak := new(alias.AesKey)
	ak.Read(stInfoKey)

	aesEnc := func(input []byte) []byte {
		enc := make([]byte, len(input))

		if dec {
			enc = alias.DecryptAES(input, ak)
		} else {
			enc = alias.EncryptAES(input, ak)
		}
		return enc
	}

	fields := reflect.TypeOf(*u)
	values := reflect.ValueOf(*u)

	num := values.NumField()

	for i := 0; i < num; i++ {
		var input string
		value := values.Field(i)
		field := fields.Field(i)

		switch value.Kind() {
		case reflect.String:
			input = value.String()
			enc := string(aesEnc([]byte(input)))
			reflect.ValueOf(u).Elem().FieldByName(field.Name).SetString(enc)
		case reflect.Int:
			continue
		}
	}
}

// HashPass method encrypt password of user
func (u *User) HashPass() {
	ak := new(alias.AesKey)
	ak.Read(passKey)

	input := []byte(u.Pass)
	encrypted := alias.EncryptAES(input, ak)

	u.Pass = alias.HashPass(string(encrypted))
}

// UpdateUser fills empty fields of user entry in DB
func (m *MysqlUser) UpdateUser(user *User) int64 {
	db := m.connect()
	defer db.Close()
	stmt := prepare(db, "UPDATE users SET mail=?, name=?, surname=?, partonymic=?, recourse=?, login=?, pass=?, activationKey=?, activationType = ?, department = ? WHERE id = ?")
	defer stmt.Close()
	res := exec(stmt, []interface{}{user.Mail, user.Name, user.Surname, user.Partonymic,
		user.Recourse, user.Login, user.Pass, user.ActivationKey, user.ActivationType, user.Department, user.ID})
	aff := affect(res)
	return aff
}

// InsertUser inserts user obj in db
func (m *MysqlUser) InsertUser(user *User) sql.Result {
	db := m.connect()
	defer db.Close()
	stmt := prepare(db, "INSERT INTO users (login, pass, mail, name, surname, partonymic, recourse, role, department, status, activationKey, activationType) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	res := exec(stmt, []interface{}{user.Login, user.Pass, user.Mail, user.Name, user.Surname, user.Partonymic, user.Recourse,
		user.Role, user.Department, user.Status, user.ActivationKey, user.ActivationType})
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
		&ns.Partonymic, &ns.Recourse, &user.Role, &user.Department, &ns.Status, &ns.ActivationKey, &ns.ActivationType)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return user
	} else if err != nil {
		panic("db error: " + err.Error())
	}
	user.chkNullFields(ns)
	return user
}

// GetUsers returns all users from db
func (m *MysqlUser) GetUsers(sqlKey string, keyVal interface{}) []*User {
	db := m.connect()
	defer db.Close()

	var rows *sql.Rows
	if sqlKey == "*" {
		rows = getQuery(db, "SELECT * FROM users")
	} else {
		stmt := prepare(db, "SELECT * FROM users WHERE "+sqlKey+" = ?")
		defer stmt.Close()
		rows = chk(stmt.Query(keyVal)).(*sql.Rows)
	}

	users := make([]*User, 0)
	for rows.Next() {
		user, ns := new(User), new(userNullFields)
		err := rows.Scan(&user.ID, &user.Login, &user.Pass, &ns.Mail, &ns.Name, &ns.Surname,
			&ns.Partonymic, &ns.Recourse, &user.Role, &user.Department, &ns.Status, &ns.ActivationKey, &ns.ActivationType)
		if err != nil {
			fmt.Println(err.Error())
		}
		user.chkNullFields(ns)
		users = append(users, user)
	}
	return users
}

// FindUser returns finded user-obj arr using key
func (m *MysqlUser) FindUser(keys []string) []*User {
	modified := false
	sqlReq := "SELECT * FROM users WHERE "

	checkModified := func() {
		if modified {
			sqlReq += "AND "
		}
		modified = true
	}

	delReservKey := func(i *int) {
		fmt.Println("DELETE: " + keys[*i])
		if *i+1 >= len(keys) {
			keys = keys[:*i]
		} else {
			keys = append(keys[:*i], keys[*i+1:]...)
		}
		*i--
	}

	for i := 0; i < len(keys); i++ {
		switch keys[i] {
		case "@admin":
			checkModified()
			delReservKey(&i)
			sqlReq += "role = \"admin\" "
		case "@user":
			checkModified()
			delReservKey(&i)
			sqlReq += "role = \"user\" "
		case "@ticketModerator":
			checkModified()
			delReservKey(&i)
			sqlReq += "role = \"ticketModerator\" "
		case "@activated":
			checkModified()
			delReservKey(&i)
			sqlReq += "name != \"\" "
		case "@inactive":
			checkModified()
			delReservKey(&i)
			sqlReq += "name = \"\" "
		}
	}

	var argsLen = 6
	addSQLReq := "(name = ? OR surname = ? OR partonymic = ? OR login = ? OR mail = ? OR recourse = ?)"
	keysLen := len(keys)
	sqlKeys := make([]string, 0)

	for i := 0; i < keysLen; i++ {
		checkModified()
		sqlReq += addSQLReq
		for j := 0; j < argsLen; j++ {
			sqlKeys = append(sqlKeys, keys[i])
		}
	}

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
			&ns.Partonymic, &ns.Recourse, &user.Role, &user.Department, &ns.Status, &ns.ActivationKey, &ns.ActivationType)
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
	if ns.ActivationType != nil {
		if reflect.TypeOf(ns.ActivationType).String() == "[]uint8" {
			u.ActivationType = chk(alias.STI(string(ns.ActivationType.([]byte)))).(int)
		} else {
			u.ActivationType = int(ns.ActivationType.(int64))
		}
	}
}

// GetDepartment method returns the department obj of user
func (u *User) GetDepartment() *Department {
	mysql := CreateMysqlUser()
	return mysql.GetDepartment("id", u.Department)
}

// HidePrivateInfo clears all private info from user obj
func (u *User) HidePrivateInfo() {
	u.Pass = ""
	//u.ID = 0
}

// String returns user object description
func (u *User) String() string {
	return "Name: " + u.Name + "\nSurname: " + u.Surname + "\nPartonymic: " + u.Partonymic +
		"\nLogin: " + u.Login + "\nMail: " + u.Mail + "\nRecourse: " + u.Recourse
}

// ChkPass returns true if regexp match the password
func (u *User) ChkPass() bool {
	if alias.StrLen(u.Pass) > 5 && alias.StrLen(u.Pass) < 20 {
		match, _ := regexp.MatchString("^[a-zA-Z0-9]{1}[a-zA-Z0-9_-]+[a-zA-Z0-9]{1}$", u.Pass)
		return match
	}
	return false
}

// ChkLogin returns true if regexp match the login
func (u *User) ChkLogin() bool {
	if alias.StrLen(u.Login) > 4 && alias.StrLen(u.Login) < 20 {
		match, _ := regexp.MatchString("^[a-zA-Z0-9]{1}[a-zA-Z0-9_-]+[a-zA-Z0-9]{1}$", u.Login)
		return match
	}
	return false
}

// DeleteUser delete`s user from db
func (m *MysqlUser) DeleteUser(sqlKey string, keyVal interface{}) int64 {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "DELETE FROM users WHERE "+sqlKey+" = ?")
	defer stmt.Close()
	res := exec(stmt, []interface{}{keyVal})
	aff := affect(res)

	return aff
}

// GetDepartmentTicketAdmin returns user obj - admin in department
func (m *MysqlUser) GetDepartmentTicketAdmin(depID int) *User {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "SELECT * FROM users WHERE department = ? AND role = ? LIMIT 1") // TODO: ticket admin
	defer stmt.Close()

	user, ns := new(User), new(userNullFields)
	err := stmt.QueryRow(depID, "ticketModerator").Scan(&user.ID, &user.Login, &user.Pass, &ns.Mail, &ns.Name, &ns.Surname,
		&ns.Partonymic, &ns.Recourse, &user.Role, &user.Department, &ns.Status, &ns.ActivationKey, &ns.ActivationType)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return user
	} else if err != nil {
		panic("db error: " + err.Error())
	}
	user.chkNullFields(ns)
	return user

}
