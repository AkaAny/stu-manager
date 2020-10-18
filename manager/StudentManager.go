package manager

import (
	"errors"
	"fmt"
	"stu-manager/model"
	"stu-manager/orm"
)

type StudentManager struct {
}

var sStudentManager *StudentManager = nil

func GetStudentManager() *StudentManager {
	if sStudentManager == nil {
		sStudentManager = &StudentManager{}
	}
	return sStudentManager
}

func (manager StudentManager) GetStudentByID(stuID int64) (model.Student, error) {
	var foundCount int64
	var stuForQuery = model.Student{ID: stuID}
	err := orm.GetStuGOrm().GetDB().First(&stuForQuery).Count(&foundCount).Error
	if err != nil {
		return stuForQuery, err
	}
	if foundCount <= 0 { //学生不存在
		var errMsg = fmt.Sprintf("student:%d is not found", stuID)
		return stuForQuery, errors.New(errMsg)
	}
	return stuForQuery, nil
}
