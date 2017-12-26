package models

import (
	"database/sql/driver"
	"fmt"
)

// EndpointResponseAction type is the go counterpart of enum for action key of endpoint response
type EndpointResponseAction string

// Scan implements the sql.Scanner interface.
func (e *EndpointResponseAction) Scan(value interface{}) error {
	if value == nil {
		*e = EndpointResponseAction_NULL
		return nil
	}

	switch v := value.(type) {
	case string:
		*e = EndpointResponseAction(v)
		return nil
	case []byte:
		*e = EndpointResponseAction(v)
		return nil
	}
	return fmt.Errorf("Can't convert %T to EndpointResponseAction", value)
}

// Value implements the driver driver.Valuer interface.
func (e EndpointResponseAction) Value() (driver.Value, error) {
	if e == EndpointResponseAction_NULL {
		return nil, nil
	}
	return string(e), nil
}

const (
	EndpointResponseAction_NULL EndpointResponseAction = ""
	EndpointResponseAction_GOTO EndpointResponseAction = "GOTO"
	EndpointResponseAction_NEXT EndpointResponseAction = "NEXT"
	EndpointResponseAction_GOTO_LATER EndpointResponseAction = "GOTO_LATER"
	EndpointResponseAction_NEXT_LATER EndpointResponseAction = "NEXT_LATER"
	EndpointResponseAction_RETRY EndpointResponseAction = "RETRY"
	EndpointResponseAction_RETRY_NOW EndpointResponseAction = "RETRY_NOW"
	EndpointResponseAction_ERROR EndpointResponseAction = "ERROR"
	EndpointResponseAction_PROBLEM EndpointResponseAction = "PROBLEM"
	EndpointResponseAction_CANCELED EndpointResponseAction = "CANCELED"
)

