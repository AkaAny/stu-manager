package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	model2 "stu-manager/controller/model"
	"stu-manager/logger"
	"stu-manager/manager"
	"stu-manager/model"
	"stu-manager/orm"
	"time"
)

type RoomController struct {
}

var sRoomController *RoomController = nil

func GetRoomController() *RoomController {
	if sRoomController == nil {
		sRoomController = &RoomController{}
	}
	return sRoomController
}

func (rc *RoomController) Init(engine *gin.Engine) {
	logger.Info.Println("init RoomController")

	err := orm.GetStuGOrm().GetDB().AutoMigrate(&model.Room{})
	if err != nil {
		logger.Error.Fatalln(err)
		return
	}
	engine.Handle(http.MethodGet, "/room", rc.handleRoomQueryAll)
	engine.Handle(http.MethodGet, "/room/:room_id", rc.handleRoomQuery)
	engine.Handle(http.MethodPost, "/room/add", rc.handleRoomAdd)
	engine.Handle(http.MethodDelete, "room/delete", rc.handleRoomDelete)
	engine.Handle(http.MethodPost, "/room/arrange", rc.handleRoomArrangeQuery)
}

func (rc RoomController) handleRoomAdd(c *gin.Context) {
	var request model2.RoomAddRequest
	if !BindJsonOrError(c, &request) {
		return
	}
	var room = model.Room{
		ID:          request.RoomID,
		BuildingID:  request.BuildingID,
		CreatedDate: time.Now()}
	if !AtomicDBOperationOrError(c, func() error {
		return orm.GetStuGOrm().GetDB().Create(&room).Error
	}) {
		return
	}
	c.JSON(http.StatusOK, model2.CreatePlainResponse("ok", "room created"))
}

func (rc RoomController) handleRoomQuery(c *gin.Context) {
	valid, roomID := ParseIntFromParamOrError(c, "room_id")
	if !valid {
		return
	}
	logger.Info.Printf("query room by id:%d", roomID)
	var roomForQuery = model.Room{ID: roomID}
	if !AtomicDBOperationOrError(c, func() error {
		var err error
		roomForQuery, err = rc.QueryRoomByRoomIDLocked(roomID)
		return err
	}) {
		return
	}
	c.JSON(http.StatusOK, roomForQuery)
}

func (rc RoomController) QueryRoomByRoomIDLocked(roomID int64) (model.Room, error) {
	var result = model.Room{ID: roomID}
	err := orm.GetStuGOrm().GetDB().First(&result).Error
	if err != nil {
		return result, err
	}
	if result.BuildingID == 0 { //未查找到指定的寝室记录
		var errMsg = fmt.Sprintf("room:%d is not found", roomID)
		return result, errors.New(errMsg)
	}
	return result, nil
}

func (rc RoomController) handleRoomQueryAll(c *gin.Context) {
	logger.Info.Println("query all rooms")
	var result = make([]model.Room, 1)
	orm.GetStuGOrm().GetDB().Find(&result)
	c.JSON(http.StatusOK, result)
	return
}

func (rc RoomController) handleRoomDelete(c *gin.Context) {
	var request model2.RoomDeleteRequest
	if !BindJsonOrError(c, &request) {
		return
	}
	var roomID = request.RoomID
	logger.Info.Printf("delete room:%d", roomID)
	if !AtomicDBOperationOrError(c, func() error {
		return orm.GetStuGOrm().GetDB().Delete(&model.Room{}, roomID).Error
	}) {
		return
	}
	var desc = fmt.Sprintf("room:%d has deleted", request.RoomID)
	c.JSON(http.StatusOK, model2.CreatePlainResponse("ok", desc))
}

func (rc RoomController) handleRoomArrangeQuery(c *gin.Context) {
	var request model2.RoomArrangeRequest
	if !BindJsonOrError(c, &request) {
		return
	}
	var dt = request.Date
	logger.Info.Printf("query arrange for room:%d date:%v", request.RoomID, dt)
	var stuToArrange model.Student
	if !AtomicDBOperationOrError(c, func() error {
		var err error
		stuToArrange, err = manager.GetRoomManager().GetArrangedStudentByDate(request.RoomID, dt)
		return err
	}) {
		return
	}
	c.JSON(http.StatusOK, stuToArrange)
}
