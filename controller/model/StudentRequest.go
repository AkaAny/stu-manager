package model

type StudentAddRequest struct {
	ID           int64  `json:"id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	RoomID       int64  `json:"room_id" binding:"required"`
	ArrangeIndex int64  `json:"arrange_index"` //自动后延或指定
}

type StudentUpdateRequest struct {
	ID           int64  `json:"id" binding:"required"` //通过学号查找原学生信息
	Name         string `json:"name"`
	RoomID       int64  `json:"room_id"`
	ArrangeIndex int64  `json:"arrange_index"`
}

type StudentDeleteRequest struct {
	ID int64 `json:"id" binding:"required"`
}
