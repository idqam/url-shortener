package model

import "time"

type Analytics struct {
	ID        string     `json:"id"`
	URLID     *string    `json:"url_id"`
	IPAddress *string    `json:"ip_address"`
	UserAgent *string    `json:"user_agent"`
	Country   *string    `json:"country"`
	ClickedAt time.Time  `json:"clicked_at"`
}
