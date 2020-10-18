package model

type Student struct {
	ID           int64  `gorm:"primaryKey;column:id"`
	Name         string `gorm:"column:name"`
	RoomID       int64  `gorm:"column:room_id"`
	ArrangeIndex int64  `gorm:"column:arrange_index"` //值日轮流顺序(index)
}
