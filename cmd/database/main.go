package main

import (
  "github.com/DEEMMOONS/BaseHTTP/internal/server"
  "github.com/DEEMMOONS/BaseHTTP/internal/database"
  "github.com/go-pg/pg"
)

func main() {
 config, err := server.CreateConfig("config/config.json")
  if err != nil {
    panic(err)
  }
  db := pg.Connect(&pg.Options{
		User:     config.DB.User,
		Password: config.DB.Password,
		Database: config.DB.Database,
	})
	defer db.Close()
	if err = database.CreateSchema(db); err != nil {
		panic(err)
	}
}
