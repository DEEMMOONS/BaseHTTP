package database

import (
  "github.com/go-pg/pg/orm"
  "github.com/go-pg/pg"
)

func AddOrder(db *pg.DB, order server.Order) error {
	if _, err := db.Model(&order).Insert(); err != nil {
		return err
	}
	return nil
}


func CreateSchema(db *pg.DB) error {
  models := []interface{}{
		(*Order)(nil),
	}
	for _, model := range models {
		op := orm.CreateTableOptions{}
		err := db.DB.Model(model).CreateTable(&op)
		if err != nil {
			return err
		}
	}
	return nil
}
