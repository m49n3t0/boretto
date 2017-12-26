package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Definition struct {
	Sequence []Step `json:"sequence"`
}

var (
	_ sql.Scanner   = &Definition{}
	_ driver.Valuer = Definition{}
)

// Scan implements the sql.Scanner interface.
func (e *Definition) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case Definition:
		// Need to deep copy struct
		*e = Definition{
			Sequence: make([]Step, len(v.Sequence)),
		}
		for idx, step := range v.Sequence {
			e.Sequence[idx] = step
		}
	case []byte:
		*e = Definition{} // json.Unmarshal does not reset the struct content
		err = json.Unmarshal(v, e)
	case string:
		*e = Definition{} // json.Unmarshal does not reset the struct content
		err = json.Unmarshal([]byte(v), e)
	default:
		return fmt.Errorf("Can't convert %T to models.Definition", value)
	}
	return
}

var emptySlice = []Step{}

// Value implements the driver driver.Valuer interface.
func (e Definition) Value() (driver.Value, error) {
	if e.Sequence == nil {
		e.Sequence = emptySlice
	}
	jsonBytes, _ := json.Marshal(e)
	return string(jsonBytes), nil
}
