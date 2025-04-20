package domain

import "time"

type Reception struct {
	ID       string     `json:"id"`
	DateTime time.Time  `json:"dateTime"`
	PvzId    string     `json:"pvzId"`
	Status   string     `json:"status"`
	ClosedAt *time.Time `json:"closedAt"`
}
