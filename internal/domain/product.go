package domain

import "time"

type Product struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionId string    `json:"receptionId"`
}
