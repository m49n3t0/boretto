package models

// HttpEndpoint is the go counterpart of table http_endpoint
type HttpEndpoint struct {
	TableName struct{} `sql:"http_endpoint"`

	Id           int64  `sql:"id"`
	Type         string `sql:"type,notnull"`
	Name         string `sql:"name,notnull"`
	CreationDate string `sql:"creation_date,notnull"`
	LastUpdate   string `sql:"last_update,notnull"`
	Method       string `sql:"method,notnull"`
	Url          string `sql:"url,notnull"`
}

const (
	TblHttpEndpoint_Id           = `"endpoint_http"."id"`
	TblHttpEndpoint_Type         = `"endpoint_http"."type"`
	TblHttpEndpoint_Name         = `"endpoint_http"."name"`
	TblHttpEndpoint_CreationDate = `"endpoint_http"."creation_date"`
	TblHttpEndpoint_LastUpdate   = `"endpoint_http"."last_update"`
	TblHttpEndpoint_Method       = `"endpoint_http"."method"`
	TblHttpEndpoint_Url          = `"endpoint_http"."url"`
)
