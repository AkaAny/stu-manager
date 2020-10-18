package model

import "time"

type RoomAddRequest struct {
	RoomID     int64 `json:"room_id" binding:"required"`
	BuildingID int64 `json:"building_id" binding:"required"`
}

type RoomDeleteRequest struct {
	RoomID int64 `json:"room_id" binding:"required"`
}

type RoomArrangeRequest struct {
	RoomID int64     `json:"room_id" binding:"required"`
	Date   time.Time `json:"date" binding:"required"`
}
