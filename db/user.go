package db

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
func (u *User) WriteIn(name, surname, partonymic, recourse, login, password string) {
	u.Name, u.Surname, u.Partonymic = name, surname, partonymic
	u.Recourse = recourse
	u.Login, u.Pass = login, password
}

// UpdateUser fills empty fields of user entry in DB
func (m *MysqlUser) UpdateUser(user *User) int64 {
	db := m.connect()
	defer db.Close()
	stmt := prepare(db, "UPDATE users SET name=? surname=? partonymic=? recourse=? login=? password=? WHERE id = ?")
	defer stmt.Close()
	res := exec(stmt, []interface{}{user.Name, user.Surname, user.Partonymic, user.Partonymic,
		user.Recourse, user.Login, user.Pass, user.ID})
	aff := affect(res)
	return aff
}

// Exist method checks if user exist
func (u *User) Exist() bool {
	if u.ID == 0 {
		return false
	}
	return true
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

	user := new(User)
	ns := new(userNullFields)
	err := stmt.QueryRow(key).Scan(&user.ID, &user.Login, &user.Pass, &ns.Mail, &ns.Name,
		&ns.Surname, &ns.Partonymic, &ns.Recourse, &user.Role, &user.Department, &ns.Status)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return user
	} else if err != nil {
		panic("db error: " + err.Error())
	}
	user.chkNullFields(ns)
	return user
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
	return mysql.GetDepartment(u.Department)
}
