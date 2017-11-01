package models

// Robot is the go counterpart of table robot
type Robot struct {
	TableName struct{} `sql:"robot"`

	Id         int64      `sql:"id"`
	Function   string     `sql:"function,notnull"`
	Version    int64      `sql:"version,notnull"`
	Definition Definition `sql:"definition,notnull"`
}

const (
	TblRobot_Id         = `"robot"."id"`
	TblRobot_Function   = `"robot"."function"`
	TblRobot_Version    = `"robot"."version"`
	TblRobot_Definition = `"robot"."definition"`
)
