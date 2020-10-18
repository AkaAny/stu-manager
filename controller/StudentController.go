package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	model2 "stu-manager/controller/model"
	"stu-manager/logger"
	"stu-manager/manager"
	"stu-manager/model"
	"stu-manager/orm"
)

type StudentController struct {
}

var sStudentController *StudentController = nil

func GetStudentController() *StudentController {
	if sStudentController == nil {
		sStudentController = &StudentController{}
	}
	return sStudentController
}

func (sc StudentController) Init(engine *gin.Engine) {
	logger.Info.Println("init StudentController")

	err := orm.GetStuGOrm().GetDB().AutoMigrate(&model.Student{})
	if err != nil {
		logger.Error.Fatalln(err)
		return
	}
	engine.Handle(http.MethodGet, "/student", sc.handleStudentQueryAll)
	engine.Handle(http.MethodGet, "/student/:stu_id", sc.handleStudentQuery)
	engine.Handle(http.MethodPost, "/student/add", sc.handleStudentAdd)
	engine.Handle(http.MethodDelete, "student/delete", sc.handleStudentDelete)
	engine.Handle(http.MethodGet, "/student/:stu_id/roommate", sc.handleRoomMateQuery)
	engine.Handle(http.MethodPut, "/student/update", sc.handleStudentUpdate)
}

func (sc StudentController) handleStudentAdd(c *gin.Context) {
	var request model2.StudentAddRequest
	if !BindJsonOrError(c, &request) {
		return
	}
	var stuToAdd = model.Student{
		ID:     request.ID,
		Name:   request.Name,
		RoomID: request.RoomID}
	if !AtomicDBOperationOrError(c, func() error {
		_, err := GetRoomController().QueryRoomByRoomIDLocked(request.RoomID)
		if err != nil { //寝室房间号不存在或其它错误
			return err
		}
		arrangeIndex, err := manager.GetRoomManager().AllocArrangeIndex(request.RoomID, request.ArrangeIndex)
		if err != nil {
			return err
		}
		stuToAdd.ArrangeIndex = arrangeIndex
		return orm.GetStuGOrm().GetDB().Create(&stuToAdd).Error
	}) {
		return
	}
	var desc = fmt.Sprintf("student:%d -> %s added", stuToAdd.ID, stuToAdd.Name)
	c.JSON(http.StatusOK, model2.CreatePlainResponse("ok", desc))
}

func (sc StudentController) handleStudentQuery(c *gin.Context) {
	valid, stuID := ParseIntFromParamOrError(c, "stu_id")
	if !valid {
		return
	}
	logger.Info.Printf("query student by id:%d", stuID)
	var stuForQuery model.Student
	if !AtomicDBOperationOrError(c, func() error {
		var err error
		stuForQuery, err = manager.GetStudentManager().GetStudentByID(stuID)
		return err
	}) {
		return
	}
	c.JSON(http.StatusOK, stuForQuery)
}

func (sc StudentController) handleStudentQueryAll(c *gin.Context) {
	logger.Info.Println("query all students")
	var result = make([]model.Student, 1)
	orm.GetStuGOrm().GetDB().Find(&result)
	c.JSON(http.StatusOK, result)
	return
}

func (sc StudentController) handleStudentDelete(c *gin.Context) {
	var request model2.StudentDeleteRequest
	if !BindJsonOrError(c, &request) {
		return
	}
	var studentID = request.ID
	logger.Info.Printf("delete student:%d", studentID)
	if !AtomicDBOperationOrError(c, func() error {
		return orm.GetStuGOrm().GetDB().Delete(&model.Student{}, studentID).Error
	}) {
		return
	}
	var desc = fmt.Sprintf("student:%d has deleted", studentID)
	c.JSON(http.StatusOK, model2.CreatePlainResponse("ok", desc))
}

func (sc StudentController) handleRoomMateQuery(c *gin.Context) {
	valid, stuID := ParseIntFromParamOrError(c, "stu_id")
	if !valid {
		return
	}
	var roomMates []model.Student
	if !AtomicDBOperationOrError(c, func() error {
		//先查询对应的寝室号
		var stuForQuery = model.Student{ID: stuID}
		orm.GetStuGOrm().GetDB().First(&stuForQuery)
		var roomID = stuForQuery.RoomID
		logger.Info.Printf("student:%d -> room:%d", stuID, roomID)
		//查询室友
		var err error
		roomMates, err = manager.GetRoomManager().GetRoomMateSlice(roomID)
		return err
	}) {
		return
	}
	c.JSON(http.StatusOK, roomMates)
}

func (sc StudentController) handleStudentUpdate(c *gin.Context) {
	var request model2.StudentUpdateRequest
	if !BindJsonOrError(c, &request) {
		return
	}
	var stuID = request.ID
	if !AtomicDBOperationOrError(c, func() error {
		//根据学号查找
		stuForUpdate, err := manager.GetStudentManager().GetStudentByID(stuID)
		if err != nil {
			return err
		}

		if request.Name != "" { //需要更新姓名
			stuForUpdate.Name = request.Name
		}
		if request.RoomID != 0 { //需要更新寝室号
			stuForUpdate.RoomID = request.RoomID
		}
		if request.ArrangeIndex != 0 { //需要更新值日安排（如果寝室号修改了就按照新寝室排）
			err := manager.GetRoomManager().UpdateOrSwapWhenDuplicate(&stuForUpdate, request.ArrangeIndex)
			if err != nil {
				return err
			}
		}
		//统一再更新一次
		return orm.GetStuGOrm().GetDB().Save(&stuForUpdate).Error
	}) {
		return
	}
	var desc = fmt.Sprintf("student:%d updated", stuID)
	c.JSON(http.StatusOK, model2.CreatePlainResponse("ok", desc))
}
