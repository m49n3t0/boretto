package models

// Task is the go counterpart of table task
type Task struct {
	TableName struct{} `sql:"task"`

	Id           int64                  `sql:"id"`
	Version      int64                  `sql:"version,notnull"`
	Context      string                 `sql:"context,notnull"`
	Function     string                 `sql:"function,notnull"`
	Step         string                 `sql:"step,notnull"`
	Status       string                 `sql:"status,notnull"`
	Retry        int64                  `sql:"retry,notnull"`
	Arguments    map[string]interface{} `sql:"arguments,notnull"`
	Buffer       map[string]interface{} `sql:"buffer,notnull"`
	CreationDate string                 `sql:"creation_date,notnull"`
	LastUpdate   string                 `sql:"last_update,notnull"`
}

const (
	TblTask_Id           = `"task"."id"`
	TblTask_Version      = `"task"."version"`
	TblTask_Context      = `"task"."context"`
	TblTask_Function     = `"task"."function"`
	TblTask_Step         = `"task"."step"`
	TblTask_Status       = `"task"."status"`
	TblTask_Retry        = `"task"."retry"`
	TblTask_Arguments    = `"task"."arguments"`
	TblTask_Buffer       = `"task"."buffer"`
	TblTask_CreationDate = `"task"."creation_date"`
	TblTask_LastUpdate   = `"task"."last_update"`
)
