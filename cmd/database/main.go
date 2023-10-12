package main



func main() {
 config, err := CreateConfig(cfgPath)
  if err != nil {
    panic(err)
  }
  db, err := pg.Connect(&pg.Options{
		User:     config.DB.User,
		Password: config.DB.Password,
		Database: config.DB.Database,
	})
	db.Open()
	defer db.Close()
	if err = database.CreateSchema(db); err != nil {
		panic(err)
	}
}
