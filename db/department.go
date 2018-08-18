package db

// Department obj
type Department struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetDepartment returns Department obj using id
func (m *MysqlUser) GetDepartment(id int) *Department {
	db := m.connect()
	defer db.Close()

	stmt := prepare(db, "SELECT * FROM departments WHERE id = ?")
	defer stmt.Close()

	d := new(Department)
	err := stmt.QueryRow(id).Scan(&d.ID, &d.Name)
	if err != nil {
		panic("db error: " + err.Error())
	}
	return d
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
