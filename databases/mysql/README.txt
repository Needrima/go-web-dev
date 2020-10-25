SQL command queries
=====================================
Connecting to workbench: db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/database_name")
=====================================
checking connection sucess: err = db.Ping()
=====================================
Selecting single row from db: query := `select * from signup where Username = ?;`
                              row := db.QueryRow(query, data)
=====================================
Selecting multiple rows from db: query := `select * from signup where;`
                                 rows, err := db.Query(query)
=====================================
passing_data_to db: stmt, err := db.Prepare(`insert signup set Firstname=?, Lastname=?, Username=?, Email=?, Password=?, Phone=?;`)
                    n, err = stmt.Exec(data, data, data, data, data, data)
=====================================
updating data in db: stmt, err := db.Prepare(`update signup set Firstname=? where Username=?;`)
                    n, err = stmt.Exec(data, data)
=====================================
delete data in db: stmt, err := db.Prepare(`delete from signup where Username=?;`)
                    n, err = stmt.Exec(data)
=====================================