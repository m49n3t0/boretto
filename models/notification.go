package models

type Notification struct {
	Action string `json:"action"`
	Data   struct {
		TaskID int64 `json:"task_id"`
	} `json:"data"`
}
