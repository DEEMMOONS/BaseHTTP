package main

import {
  "github.com/DEEMMOONS/BaseHTTP/internal/server"
}

func main() {
 config, err := server.CreateConfig(cfgPath)
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
