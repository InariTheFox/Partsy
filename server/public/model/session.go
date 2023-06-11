package model

type Session struct {
	Id             string `json:"id"`
	Token          string `json:"token"`
	CreatedAt      int64  `json:"created_at"`
	ExpiresAt      int64  `json:"expires_at"`
	LastActivityAt int64  `json:"last_activity_at"`
	UserId         string `json:"user_id"`
	Roles          string `json:"roles"`
	IsOAuth        bool   `json:"is_oauth"`
}
