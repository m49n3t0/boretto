package models

//import (
//	"database/sql"
//	"database/sql/driver"
//	"encoding/json"
//	"fmt"
//)
//
//type EndpointResponse struct {
//	Action EndpointResponseAction `json:"action,notnull"`
//	Buffer *JsonB                 `json:"buffer"`
//	Data   *EndpointResponseData  `json:"data"`
//}

//var (
//	_ sql.Scanner   = &EndpointResponse{}
//	_ driver.Valuer = EndpointResponse{}
//)
//
//// Scan implements the sql.Scanner interface.
//func (e *EndpointResponse) Scan(value interface{}) (err error) {
//	switch v := value.(type) {
//	case EndpointResponse:
//		// Need to deep copy struct
//		*e = EndpointResponse{
//			Sequence: make([]Step, len(v.Sequence)),
//		}
//		for idx, step := range v.Sequence {
//			e.Sequence[idx] = step
//		}
//	case []byte:
//		*e = EndpointResponse{} // json.Unmarshal does not reset the struct content
//		err = json.Unmarshal(v, e)
//	case string:
//		*e = EndpointResponse{} // json.Unmarshal does not reset the struct content
//		err = json.Unmarshal([]byte(v), e)
//	default:
//		return fmt.Errorf("Can't convert %T to models.EndpointResponse", value)
//	}
//	return
//}
//
//var emptySlice = []Step{}
//
//// Value implements the driver driver.Valuer interface.
//func (e EndpointResponse) Value() (driver.Value, error) {
//	if e.Sequence == nil {
//		e.Sequence = emptySlice
//	}
//	jsonBytes, _ := json.Marshal(e)
//	return string(jsonBytes), nil
//}
