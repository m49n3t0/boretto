package models

import (
	"time"
)

// Machine is the go counterpart of table machine
type Machine struct {
	TableName struct{} `sql:"machine"`

	Id           int64      `sql:"id"`
	Function     string     `sql:"function,notnull"`
	Version      int64      `sql:"version,notnull"`
	Status       bool       `sql:"status,notnull"`
	Definition   Definition `sql:"definition,notnull"`
	CreationDate time.Time  `sql:"creation_date,notnull"`
	LastUpdate   time.Time  `sql:"last_update,notnull"`
}

const (
	TblMachine_Id           = `"machine"."id"`
	TblMachine_Function     = `"machine"."function"`
	TblMachine_Version      = `"machine"."version"`
	TblMachine_Status       = `"machine"."status"`
	TblMachine_Definition   = `"machine"."definition"`
	TblMachine_CreationDate = `"machine"."creation_date"`
	TblMachine_LastUpdate   = `"machine"."last_update"`
)

const (
	ColMachine_Id           = `"id"`
	ColMachine_Function     = `"function"`
	ColMachine_Version      = `"version"`
	ColMachine_Status       = `"status"`
	ColMachine_Definition   = `"definition"`
	ColMachine_CreationDate = `"creation_date"`
	ColMachine_LastUpdate   = `"last_update"`
)
