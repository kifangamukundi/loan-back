package models

type GroupMember struct {
	GroupID  uint `gorm:"primaryKey"`
	MemberID uint `gorm:"primaryKey"`
}
