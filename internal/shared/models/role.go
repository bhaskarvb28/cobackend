package models

type Role struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	DisplayName string `json:"display_name"`
}