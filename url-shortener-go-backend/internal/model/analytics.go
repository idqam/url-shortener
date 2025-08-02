package model

import "time"

type Analytics struct {
	ID         string    `json:"id"`
	URLID      *string   `json:"url_id"`
	Country    *string   `json:"country"`
	Referrer   *string   `json:"referrer"`
	UserAgent  *string   `json:"user_agent"`
	DeviceType *string   `json:"device_type"`
	Browser    *string   `json:"browser"`
	OS         *string   `json:"os"`
	ClickedAt  time.Time `json:"clicked_at"`
}
