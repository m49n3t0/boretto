package models

import (
	"database/sql/driver"
	"fmt"
)

// TaskStatus type is the go counterpart of DB enum task_status
type TaskStatus string

// Scan implements the sql.Scanner interface.
func (e *TaskStatus) Scan(value interface{}) error {
	if value == nil {
		*e = TaskStatus_NULL
		return nil
	}

	switch v := value.(type) {
	case string:
		*e = TaskStatus(v)
		return nil
	case []byte:
		*e = TaskStatus(v)
		return nil
	}
	return fmt.Errorf("Can't convert %T to TaskStatus", value)
}

// Value implements the driver driver.Valuer interface.
func (e TaskStatus) Value() (driver.Value, error) {
	if e == TaskStatus_NULL {
		return nil, nil
	}
	return string(e), nil
}

const (
	TaskStatus_NULL      TaskStatus = ""
	TaskStatus_TODO      TaskStatus = "TODO"
	TaskStatus_DOING     TaskStatus = "DOING"
	TaskStatus_ERROR     TaskStatus = "ERROR"
	TaskStatus_PROBLEM   TaskStatus = "PROBLEM"
	TaskStatus_DONE      TaskStatus = "DONE"
	TaskStatus_CANCELED TaskStatus = "CANCELED"
)




