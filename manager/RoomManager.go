package manager

import (
	"errors"
	"fmt"
	"stu-manager/logger"
	"stu-manager/model"
	"stu-manager/orm"
	"stu-manager/utils"
	"time"
)

type RoomManager struct {
}

var sRoomManager *RoomManager = nil

//单例模式很想用IOC，但BeanDefinition不知道能不能用go实现
func GetRoomManager() *RoomManager {
	if sRoomManager == nil {
		sRoomManager = &RoomManager{}
	}
	return sRoomManager
}

func (manager RoomManager) GetRoomByRoomID(roomID int64) (model.Room, error) {
	var db = orm.GetStuGOrm().GetDB()
	var foundCount int64
	var roomForQuery = model.Room{ID: roomID}
	err := db.First(&roomForQuery).Count(&foundCount).Error
	if err != nil {
		return roomForQuery, err
	}
	if foundCount <= 0 {
		var errMsg = fmt.Sprintf("room:%d is not found", roomID)
		return roomForQuery, errors.New(errMsg)
	}
	return roomForQuery, nil
}

func (manager RoomManager) GetRoomMateSlice(roomID int64) ([]model.Student, error) {
	var db = orm.GetStuGOrm().GetDB()
	var roomMateCount int64
	var roomMates = make([]model.Student, 1)
	err := db.Where("room_id", roomID).Find(&roomMates).Count(&roomMateCount).Error
	if roomMateCount <= 0 {
		return nil, err
	}
	return roomMates, err
}

func (manager RoomManager) AllocArrangeIndex(roomID int64, arrangeIndex int64) (int64, error) {
	roomMates, err := manager.GetRoomMateSlice(roomID)
	if err != nil {
		return 0, nil
	}
	//var defaultArrangeIndex=int64(len(roomMates))+1 //自动分配到最后一个
	//获取最小的不重复index，再次吐槽go没泛型，C#直接list.Select(()=>{return stu.ArrangeIndex;}就完事了
	var arrangeIndexes []int64
	for _, stu := range roomMates {
		arrangeIndexes = append(arrangeIndexes, stu.ArrangeIndex)
	}
	defaultArrangeIndex := utils.GetMinInsertValueWhenGap(arrangeIndexes, 1)
	if arrangeIndex == 0 { //没有手动指定arrangeIndex
		return defaultArrangeIndex, nil
	}
	if arrangeIndex > defaultArrangeIndex {
		logger.Warning.Printf("arrange index:%d cannot be more than current roommate count:%d",
			arrangeIndex, defaultArrangeIndex)
		return defaultArrangeIndex, nil
	}
	for _, stu := range roomMates {
		if stu.ArrangeIndex == arrangeIndex { //arrangeIndex已经被分配
			logger.Warning.Printf("arrange index:%d has been allocated,use default:%d",
				arrangeIndex, defaultArrangeIndex)
			return defaultArrangeIndex, nil
		}
	}
	return arrangeIndex, nil
}

// 注意这里的stuToUpdate必须是引用传递，因为会更新字段值
func (manager RoomManager) UpdateOrSwapWhenDuplicate(stuToUpdate *model.Student, arrangeIndexToUpdate int64) error {
	roomMates, err := manager.GetRoomMateSlice(stuToUpdate.RoomID)
	if err != nil {
		return err
	}
	var roomMateCount = int64(len(roomMates))
	if roomMateCount < arrangeIndexToUpdate {
		var errMsg = fmt.Sprintf("%d cannot bigger than room mate count:%d",
			arrangeIndexToUpdate, roomMateCount)
		return errors.New(errMsg)
	}
	var stuToSwap model.Student
	for _, stu := range roomMates { //获取重复arrange index的室友
		if arrangeIndexToUpdate != stu.ArrangeIndex {
			continue
		}
		stuToSwap = stu
	}
	if stuToSwap.ID == 0 { //没有重复
		stuToUpdate.ArrangeIndex = arrangeIndexToUpdate
		orm.GetStuGOrm().GetDB().Save(stuToUpdate)
		return nil
	} else { //arrange index重复，与当前室友互换
		stuToSwap.ArrangeIndex, stuToUpdate.ArrangeIndex = stuToUpdate.ArrangeIndex, stuToSwap.ArrangeIndex
		orm.GetStuGOrm().GetDB().Save(stuToUpdate)
		orm.GetStuGOrm().GetDB().Save(&stuToSwap)
	}
	return nil
}

func (manager RoomManager) GetArrangedStudentByDate(roomID int64, dt time.Time) (model.Student, error) {
	var stuToArrange model.Student
	room, err := manager.GetRoomByRoomID(roomID)
	if err != nil {
		return stuToArrange, err
	}
	roomMates, err := manager.GetRoomMateSlice(roomID)
	if err != nil {
		return stuToArrange, err
	}
	var roomMateCount = int64(len(roomMates))
	//计算天数差
	var offsetInDay = dt.Sub(room.CreatedDate) / (time.Hour * 24)
	//取余数+1得到当天的值日index
	var arrangeIndex = (int64(offsetInDay) % roomMateCount) + 1
	logger.Info.Printf("day offset:%d -> arrange index:%d", offsetInDay, arrangeIndex)
	//查找arrangeIndex对应的室友
	for _, stu := range roomMates {
		if arrangeIndex == stu.ArrangeIndex {
			stuToArrange = stu
			break
		}
	}
	if stuToArrange.ID == 0 { //未找到arrangeIndex对应的室友
		var errMsg = fmt.Sprintf("student with index:%d is not found in room:%d", arrangeIndex, roomID)
		return stuToArrange, errors.New(errMsg)
	}
	return stuToArrange, nil
}
