package models

// HttpEndpoint is the go counterpart of table http_endpoint
type HttpEndpoint struct {
	TableName struct{} `sql:"http_endpoint"`

	Id           int64  `sql:"id"`
	Type         string `sql:"type,notnull"`
	Version      int64  `sql:"version,notnull"`
	Name         string `sql:"name,notnull"`
	CreationDate string `sql:"creation_date,notnull"`
	LastUpdate   string `sql:"last_update,notnull"`
	Method       string `sql:"method,notnull"`
	Url          string `sql:"url,notnull"`
}

const (
	TblHttpEndpoint_Id           = `"http_endpoint"."id"`
	TblHttpEndpoint_Type         = `"http_endpoint"."type"`
	TblHttpEndpoint_Version      = `"http_endpoint"."version"`
	TblHttpEndpoint_Name         = `"http_endpoint"."name"`
	TblHttpEndpoint_CreationDate = `"http_endpoint"."creation_date"`
	TblHttpEndpoint_LastUpdate   = `"http_endpoint"."last_update"`
	TblHttpEndpoint_Method       = `"http_endpoint"."method"`
	TblHttpEndpoint_Url          = `"http_endpoint"."url"`
)
