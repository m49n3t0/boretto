package models

import (
	"time"
)

// Endpoint is the go counterpart of table endpoint
type Endpoint struct {
	TableName struct{} `sql:"endpoint"`

	Id           int64     `sql:"id"`
	Version      int64     `sql:"version,notnull"`
	Name         string    `sql:"name,notnull"`
	Method       string    `sql:"method,notnull"`
	Url          string    `sql:"url,notnull"`
	CreationDate time.Time `sql:"creation_date,notnull"`
	LastUpdate   time.Time `sql:"last_update,notnull"`
}

const (
	TblEndpoint_Id           = `"endpoint"."id"`
	TblEndpoint_Version      = `"endpoint"."version"`
	TblEndpoint_Name         = `"endpoint"."name"`
	TblEndpoint_Method       = `"endpoint"."method"`
	TblEndpoint_Url          = `"endpoint"."url"`
	TblEndpoint_CreationDate = `"endpoint"."creation_date"`
	TblEndpoint_LastUpdate   = `"endpoint"."last_update"`
)

const (
	ColEndpoint_Id           = `"id"`
	ColEndpoint_Version      = `"version"`
	ColEndpoint_Name         = `"name"`
	ColEndpoint_Method       = `"method"`
	ColEndpoint_Url          = `"url"`
	ColEndpoint_CreationDate = `"creation_date"`
	ColEndpoint_LastUpdate   = `"last_update"`
)
