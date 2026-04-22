package models

type SystemSetting struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PeriodStatusResponse struct {
	IsOpen bool `json:"is_open"`
}

type PeriodUpdateRequest struct {
	IsOpen bool `json:"is_open" binding:"required"`
}
