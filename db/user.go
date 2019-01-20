package db

import (
	"CypherDesk-main/alias"
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

var (
	findUserFields = []string{"Login", "Mail", "Name", "Surname", "Partonymic", "Recourse"}
)

const (
	StInfoKey        = "keys/userdatakey.toml"
	PassKey          = "keys/passkey.toml"
	ActivationKeyKey = "keys/activationKey.toml"
	IDKey            = "keys/idkey.toml"
)

// User - struct describing the registered user
type User struct {
	ID         int    `json: "id"`
	Login      string `json: "login"`
	Pass       string `json: "pass"`
	Mail       string `json: "mail"`
	Name       string `json: "name"`
	Surname    string `json: "surname"`
	Partonymic string `json: "partonymic"`
	Recourse   string `json: "recourse"`
	Role       string `json: "role"`
	Department int    `json: "department"`
	Status     string `json: "status"`
}

type userNullFields struct {
	Mail       interface{}
	Name       interface{}
	Surname    interface{}
	Partonymic interface{}
	Recourse   interface{}
	Status     interface{}
}

func DecID(encID string) int {
	decID, err := alias.STI(alias.StandartRefact(alias.Base32Decode(encID), true, IDKey))
	if err != nil {
		fmt.Println("Error when decrypt uder ID")
	}
	return decID
}

func (u *User) GetEncID() string {
	strID := strconv.Itoa(u.ID)
	return alias.Base32Encode(alias.StandartRefact(strID, false, IDKey))
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

// WriteIn fills empty fields fo user obj
func (u *User) WriteIn(user *User) {
	u.Name, u.Surname, u.Partonymic = user.Name, user.Surname, user.Partonymic
	u.Recourse = user.Recourse
	u.Login, u.Pass = user.Login, user.Pass

	if !alias.EmptyStr(user.Mail) {
		u.Mail = user.Mail
	}
}

// RefactField method decrypts or encrypts user object field
func (u *User) RefactField(fieldName string, dec bool) {
	ak := new(alias.AesKey)
	ak.Read(StInfoKey)

	r := reflect.ValueOf(*u)
	f := reflect.Indirect(r).FieldByName(fieldName)

	var input string

	switch f.Kind() {
	case reflect.String:
		input = f.String()
		enc := alias.StandartRefact(input, dec, StInfoKey)
		reflect.ValueOf(u).Elem().FieldByName(fieldName).SetString(enc)
	}
}

// RefactStandartInfo encrypts/decrypts all string-fields of user
func (u *User) RefactStandartInfo(dec bool) {
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

			var enc string
			if field.Name == "Pass" {
				enc = alias.StandartRefact(input, dec, PassKey)
			} else {
				enc = alias.StandartRefact(input, dec, StInfoKey)
			}

			reflect.ValueOf(u).Elem().FieldByName(field.Name).SetString(enc)
		case reflect.Int:
			continue
		}
	}
}

// HashPass method encrypt password of user
func (u *User) HashPass() {
	u.Pass = alias.HashPass(u.Pass)
}

// UpdateUser fills empty fields of user entry in DB
func (m *MysqlUser) UpdateUser(user *User) int64 {
	user.RefactStandartInfo(false)
	db := m.connect()
	defer db.Close()
	stmt := prepare(db, "UPDATE users SET mail=?, name=?, surname=?, partonymic=?, recourse=?, login=?, pass=?, department = ? WHERE id = ?")
	defer stmt.Close()
	res := exec(stmt, []interface{}{user.Mail, user.Name, user.Surname, user.Partonymic,
		user.Recourse, user.Login, user.Pass, user.Department, user.ID})
	aff := affect(res)
	return aff
}

// InsertUser inserts user obj in db
func (m *MysqlUser) InsertUser(user *User) sql.Result {
	user.RefactStandartInfo(false)
	db := m.connect()
	defer db.Close()
	stmt := prepare(db, "INSERT INTO users (login, pass, mail, name, surname, partonymic, recourse, role, department, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	res := exec(stmt, []interface{}{user.Login, user.Pass, user.Mail, user.Name, user.Surname, user.Partonymic, user.Recourse,
		user.Role, user.Department, user.Status})
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

// GetUserByDecField finds user by decrypted val
func (m *MysqlUser) GetUserByDecField(sqlParam string, key interface{}) *User {
	var refactKey string
	switch v := key.(type) {
	case string:
		switch key {
		case "Pass":
			v = alias.StandartRefact(v, false, PassKey)
		default:
			v = alias.StandartRefact(v, false, StInfoKey)
		}
		refactKey = v
	}
	return m.GetUser(sqlParam, refactKey)
}

// GetUser returns user obj using diff keys (id, login, mail)
func (m *MysqlUser) GetUser(sqlParam string, key interface{}) *User {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "SELECT * FROM users WHERE "+sqlParam+" = ? LIMIT 1")
	defer stmt.Close()

	user, ns := new(User), new(userNullFields)
	err := stmt.QueryRow(key).Scan(&user.ID, &user.Login, &user.Pass, &ns.Mail, &ns.Name, &ns.Surname,
		&ns.Partonymic, &ns.Recourse, &user.Role, &user.Department, &ns.Status)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return user
	} else if err != nil {
		panic("db error: " + err.Error())
	}
	user.chkNullFields(ns)
	user.RefactStandartInfo(true)
	return user
}

func (m *MysqlUser) GetUsersByDecField(sqlKey string, keyVal interface{}) []*User {
	switch t := keyVal.(type) {
	case string:
		return m.GetUsers(sqlKey, alias.StandartRefact(t, false, StInfoKey))
	}
	return m.GetUsers(sqlKey, keyVal)
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
			&ns.Partonymic, &ns.Recourse, &user.Role, &user.Department, &ns.Status)
		if err != nil {
			fmt.Println(err.Error())
		}
		user.chkNullFields(ns)
		user.RefactStandartInfo(true)
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
			sqlReq += "role = \""+EncryptedAdminValue+"\" "
		case "@user":
			checkModified()
			delReservKey(&i)
			sqlReq += "role = \""+EncryptedUserValue+"\" "
		case "@ticketModerator":
			checkModified()
			delReservKey(&i)
			sqlReq += "role = \""+EncryptedTicketModeratorValue+"\" "
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
		InterfaceArgs[i] = alias.StandartRefact(a, false, StInfoKey)
	}

	rows := chk(stmt.Query(InterfaceArgs...)).(*sql.Rows)

	for rows.Next() {
		user, ns := new(User), new(userNullFields)
		err := rows.Scan(&user.ID, &user.Login, &user.Pass, &ns.Mail, &ns.Name, &ns.Surname,
			&ns.Partonymic, &ns.Recourse, &user.Role, &user.Department, &ns.Status)
		if err != nil {
			fmt.Println(err.Error())
		}
		user.chkNullFields(ns)
		user.RefactStandartInfo(true)
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
	if ns.Status != nil {
		u.Status = string(ns.Status.([]byte))
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
	err := stmt.QueryRow(depID, EncryptedTicketModeratorValue).Scan(&user.ID, &user.Login, &user.Pass, &ns.Mail, &ns.Name, &ns.Surname,
		&ns.Partonymic, &ns.Recourse, &user.Role, &user.Department, &ns.Status)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return user
	} else if err != nil {
		panic("db error: " + err.Error())
	}
	user.chkNullFields(ns)
	user.RefactStandartInfo(true)
	return user
}
