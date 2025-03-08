package models

type ComputerStatus string

const (
	Free ComputerStatus = "free"
	Busy ComputerStatus = "busy"
)

type Computer struct {
	ID       int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	PCNumber int            `gorm:"uniqueIndex" json:"pc_number"`
	Status   ComputerStatus `json:"status"`
}
