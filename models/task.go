package models

import (
	"time"
)

// Task is the go counterpart of table task
type Task struct {
	TableName struct{} `sql:"task"`

	Id           int64                  `sql:"id"`
	Version      int64                  `sql:"version,notnull"`
	Context      string                 `sql:"context,notnull"`
	Function     string                 `sql:"function,notnull"`
	Step         string                 `sql:"step,notnull"`
	Status       TaskStatus                 `sql:"status,notnull"`
	Retry        int64                  `sql:"retry,notnull"`
	CreationDate time.Time              `sql:"creation_date,notnull"`
	LastUpdate   time.Time              `sql:"last_update,notnull"`
	TodoDate     time.Time              `sql:"todo_date,notnull"`
	DoneDate     *time.Time             `sql:"done_date"`
	Arguments    map[string]interface{} `sql:"arguments,notnull"`
	Buffer       map[string]interface{} `sql:"buffer,notnull"`
    Comment      string                 `sql:"comment"`
}

const (
	TblTask_Id           = `"task"."id"`
	TblTask_Version      = `"task"."version"`
	TblTask_Context      = `"task"."context"`
	TblTask_Function     = `"task"."function"`
	TblTask_Step         = `"task"."step"`
	TblTask_Status       = `"task"."status"`
	TblTask_Retry        = `"task"."retry"`
	TblTask_CreationDate = `"task"."creation_date"`
	TblTask_LastUpdate   = `"task"."last_update"`
	TblTask_TodoDate     = `"task"."todo_date"`
	TblTask_DoneDate     = `"task"."done_date"`
	TblTask_Arguments    = `"task"."arguments"`
	TblTask_Buffer       = `"task"."buffer"`
    TblTask_Comment     = `"task"."comment"`
)

const (
	ColTask_Id           = `"id"`
	ColTask_Version      = `"version"`
	ColTask_Context      = `"context"`
	ColTask_Function     = `"function"`
	ColTask_Step         = `"step"`
	ColTask_Status       = `"status"`
	ColTask_Retry        = `"retry"`
	ColTask_CreationDate = `"creation_date"`
	ColTask_LastUpdate   = `"last_update"`
	ColTask_TodoDate     = `"todo_date"`
	ColTask_DoneDate     = `"done_date"`
	ColTask_Arguments    = `"arguments"`
	ColTask_Buffer       = `"buffer"`
    ColTask_Comment      = `"comment"`
)
