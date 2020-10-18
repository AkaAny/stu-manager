package orm

import (
	"gorm.io/gorm"
	"stu-manager/logger"
)

type Creator interface {
	Create(dbConfig DBConfig) (*gorm.DB, error)
}

var sSelectorMap = make(map[string]Creator)

func GetCreatorByType(dbType string) Creator {
	return sSelectorMap[dbType]
}

func init() {
	logger.Info.Println("init db creators")
	sSelectorMap["mysql"] = MySQLCreator{}
}
