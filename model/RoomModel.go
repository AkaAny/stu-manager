package model

import (
	"time"
)

type Room struct {
	ID          int64     `gorm:"primaryKey;column:id" json:"id"`
	CreatedDate time.Time `gorm:"column:created_date" json:"created_date"` //通过CreatedDate来计算值日开始轮流的日期
	BuildingID  int64     `gorm:"column:building_id" json:"building_id"`
}
