package db

import (
	"CypherDesk-main/alias"
	"database/sql"
	"reflect"
	"time"
)

const (
	ExtUserKey = "keys/extendedUserKey.toml"
)

type ExtendedUser struct {
	ID                  int        `json: "id"`
	ActivationKey       string     `json: "activationKey"`
	ActivationType      string     `json: "activationType"`
	ActivationDate      time.Time  `json: "activationDate"`
	Phone               string     `json: "phone"`
	AccountCreationDate time.Time  `json: "accountCreationDate"`
	Address             string     `json: "address"`
}

// RefactStandartInfo encrypts/decrypts all string-fields of user
func (eu *ExtendedUser) Refact(dec bool) {
	ak := new(alias.AesKey)
	ak.Read(ExtUserKey)

	fields := reflect.TypeOf(*eu)
	values := reflect.ValueOf(*eu)

	num := values.NumField()

	for i := 0; i < num; i++ {
		var input string
		value := values.Field(i)
		field := fields.Field(i)

		switch value.Kind() {
		case reflect.String:
			input = value.String()

			var enc string
			if input != "" && input != "none" {
				enc = alias.StandartRefact(input, dec, StInfoKey)
			}

			reflect.ValueOf(eu).Elem().FieldByName(field.Name).SetString(enc)
		case reflect.Int:
			continue
		}
	}
}

// GetExtendedUser returns extended_user object by standart_user(id)
func (m *MysqlUser) GetExtenedUser(user *User) (*ExtendedUser) {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "SELECT * FROM extended_user WHERE user_id = ? LIMIT 1")
	defer stmt.Close()

	extUser := new(ExtendedUser)
	err := stmt.QueryRow(user.ID).Scan(&extUser.ID, &extUser.ActivationKey, &extUser.ActivationType, &extUser.ActivationDate,
		&extUser.Phone, &extUser.AccountCreationDate, &extUser.Address)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return extUser
	} else if err != nil {
		panic("db error: " + err.Error())
	}
	extUser.Refact(true)
	return extUser
}

// SetActivationKey
func (eu *ExtendedUser) SetKey(keyType string) {
	eu.ActivationKey = alias.MD5(alias.StringWithCharset(20, loginCharset) + string(eu.ID) + time.Now().String())
	eu.ActivationDate = time.Now()
	eu.ActivationType = keyType
}

func (m *MysqlUser) UpdateExtendedUserField(eu *ExtendedUser, name string, value interface{}) int64 {
	eu.Refact(false)
	db := m.connect()
	defer db.Close()
	stmt := prepare(db, "UPDATE extended_user SET "+name+"=? WHERE id = ?")
	defer stmt.Close()
	res := exec(stmt, []interface{}{value, eu.ID})
	aff := affect(res)
	return aff
}

func (m *MysqlUser) InsertExtendedUser(eu *ExtendedUser) sql.Result {
	eu.Refact(false)
	db := m.connect()
	defer db.Close()
	stmt := prepare(db, "INSERT INTO extended_user (user_id, activationKey, activationType, activationDate, phone, accountCreationDate, address) VALUES (?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	res := exec(stmt, []interface{}{eu.ID, eu.ActivationKey, eu.ActivationType, eu.ActivationDate, eu.Phone, eu.AccountCreationDate, eu.Address})
	return res
}

func (m *MysqlUser) UpdateExtendedUser(eu *ExtendedUser) int64 {
	eu.Refact(false)
	db := m.connect()
	defer db.Close()
	stmt := prepare(db, "UPDATE extended_user SET activationKey=?, activationType=?, activationDate=?, phone=?, accountCreationDate=?, address=? WHERE user_id=?")
	defer stmt.Close()
	res := exec(stmt, []interface{}{eu.ActivationKey, eu.ActivationType, eu.ActivationDate, eu.Phone, eu.AccountCreationDate,
		eu.Address, eu.ID})
	aff := affect(res)
	return aff
}