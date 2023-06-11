package model

type Part struct {
	Id          string         `json:"id"`
	CreatedAt   int64          `json:"created_at"`
	CreatedBy   string         `json:"created_by"`
	UpdatedAt   int64          `json:"updated_at"`
	UpdatedBy   string         `json:"updated_by"`
	DeletedAt   int64          `json:"deleted_at"`
	DeletedBy   string         `json:"deleted_by"`
	Type        PartType       `json:"part_type"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Props       map[string]any `json:"props"`
}

type PartList []*Part

type PartType string

type PartsWithCount struct {
	Parts      PartList `json:"parts"`
	TotalCount int64    `json:"total_count"`
}
