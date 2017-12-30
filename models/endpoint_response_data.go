package models

//import (
//	"database/sql"
//	"database/sql/driver"
//	"encoding/json"
//	"fmt"
//)
//
//type EndpointResponseData struct {
//	Step        *string           `json:"step"`
//	Interval    *int64            `json:"interval"`
//	Comment     *string           `json:"comment"`
//	Detail      map[string]string `json:"detail"`
//	NoDecrement *bool             `json:"no_decrement"`
//}
//
//var (
//	_ sql.Scanner   = &EndpointResponseData{}
//	_ driver.Valuer = EndpointResponseData{}
//)
//
//
//// Scan implements the sql.Scanner interface.
//func (e *EndpointResponseData) Scan(value interface{}) (err error) {
//	switch v := value.(type) {
//	case EndpointResponseData:
//		// Need to deep copy struct
//		*e = EndpointResponseData{
//			Sequence: make([]Step, len(v.Sequence)),
//		}
//		for idx, step := range v.Sequence {
//			e.Sequence[idx] = step
//		}
//	case []byte:
//		*e = EndpointResponseData{} // json.Unmarshal does not reset the struct content
//		err = json.Unmarshal(v, e)
//	case string:
//		*e = EndpointResponseData{} // json.Unmarshal does not reset the struct content
//		err = json.Unmarshal([]byte(v), e)
//	default:
//		return fmt.Errorf("Can't convert %T to models.EndpointResponseData", value)
//	}
//	return
//}
//
//var emptySlice = []Step{}
//
//// Value implements the driver driver.Valuer interface.
//func (e EndpointResponseData) Value() (driver.Value, error) {
//	if e.Sequence == nil {
//		e.Sequence = emptySlice
//	}
//	jsonBytes, _ := json.Marshal(e)
//	return string(jsonBytes), nil
//}
