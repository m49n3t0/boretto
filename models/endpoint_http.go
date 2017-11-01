package models

// EndpointHttp is the go counterpart of table endpoint_http
type EndpointHttp struct {
	TableName struct{} `sql:"endpoint_http"`

	Id     int64  `sql:"id"`
	Name   string `sql:"name,notnull"`
	Method string `sql:"method,notnull"`
	Url    string `sql:"url,notnull"`
}

const (
	TblEndpointHttp_Id     = `"endpoint_http"."id"`
	TblEndpointHttp_Name   = `"endpoint_http"."name"`
	TblEndpointHttp_Method = `"endpoint_http"."method"`
	TblEndpointHttp_Url    = `"endpoint_http"."url"`
)
