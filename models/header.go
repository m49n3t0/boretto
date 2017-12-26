package models

import (
	"time"
)

// Header is the go counterpart of table header
type Header struct {
	TableName struct{} `sql:"header"`

	Id           int64     `sql:"id"`
	Name         string    `sql:"name,notnull"`
	Value        string    `sql:"value,notnull"`
	CreationDate time.Time `sql:"creation_date,notnull"`
	LastUpdate   time.Time `sql:"last_update,notnull"`
}

const (
	TblHeader_Id           = `"header"."id"`
	TblHeader_Name         = `"header"."name"`
	TblHeader_Value        = `"header"."value"`
	TblHeader_CreationDate = `"header"."creation_date"`
	TblHeader_LastUpdate   = `"header"."last_update"`
)

const (
	ColHeader_Id           = `"id"`
	ColHeader_Name         = `"name"`
	ColHeader_Value        = `"value"`
	ColHeader_CreationDate = `"creation_date"`
	ColHeader_LastUpdate   = `"last_update"`
)
