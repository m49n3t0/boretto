package models

// Robot is the go counterpart of table robot
type Robot struct {
	TableName struct{} `sql:"robot"`

	Id           int64      `sql:"id"`
	Function     string     `sql:"function,notnull"`
	Version      int64      `sql:"version,notnull"`
	Status       bool       `sql:"status,notnull"`
	Definition   Definition `sql:"definition,notnull"`
	CreationDate string     `sql:"creation_date,notnull"`
	LastUpdate   string     `sql:"last_update,notnull"`
}

const (
	TblRobot_Id           = `"robot"."id"`
	TblRobot_Function     = `"robot"."function"`
	TblRobot_Version      = `"robot"."version"`
	TblRobot_Status       = `"robot"."status"`
	TblRobot_Definition   = `"robot"."definition"`
	TblRobot_CreationDate = `"robot"."creation_date"`
	TblRobot_LastUpdate   = `"robot"."last_update"`
)
