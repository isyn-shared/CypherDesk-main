package db

// User - struct describing the registered user
type User struct {
	Id         int    `json: "id"`
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
	mail       *NullString
	name       *NullString
	surname    *NullString
	partonymic *NullString
	recourse   *NullString
	status     *NullString
}

// BasicUser returns user obj containing system information
func BasicUser(mail, role, status string, department int) *User {
	return &User{
		Mail:       mail,
		Role:       role,
		Status:     status,
		Department: department,
	}
}

// GetUser return user obj using diff keys (id, login, mail)
func (m *MysqlUser) GetUser(sqlParam string, key interface{}) *User {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "SELECT * FROM users WHERE "+sqlParam+" = ? LIMIT 1")
	defer stmt.Close()

	user := new(User)
	nf := new(userNullFields)
	err := stmt.QueryRow(key).Scan(&user.Id, &user.Login, &user.Pass, &nf.mail, &nf.name,
		&nf.surname, &nf.partonymic, &nf.recourse, &user.Role, &user.Department, &nf.status)
	if err != nil {
		panic("db error: " + err.Error())
	}

	user.chkUserNullFields(nf)

	return user
}

// TODO: think about refactoring
func (u *User) chkUserNullFields(nf *userNullFields) {
	if nf.mail != nil {
		u.Mail = nf.mail.String
	}
	if nf.name != nil {
		u.Name = nf.name.String
	}
	if nf.partonymic != nil {
		u.Partonymic = nf.partonymic.String
	}
	if nf.recourse != nil {
		u.Recourse = nf.recourse.String
	}
	if nf.status != nil {
		u.Status = nf.status.String
	}
	if nf.surname != nil {
		u.Surname = nf.surname.String
	}
}
