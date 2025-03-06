package models

import "gorm.io/gorm"

type ComputerStatus string

const (
	Free ComputerStatus = "free"
	Busy ComputerStatus = "busy"
)

type Computer struct {
	gorm.Model
	PCNumber int            `gorm:"uniqueIndex" json:"pc_number"`
	Status   ComputerStatus `json:"status"`
}
