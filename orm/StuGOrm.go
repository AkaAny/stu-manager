package orm

import (
	"github.com/pelletier/go-toml"
	"gorm.io/gorm"
	"stu-manager/config"
	"stu-manager/logger"
)

type StuGOrm struct {
	mDB *gorm.DB
}

var sStuGOrm *StuGOrm = nil

func GetStuGOrm() *StuGOrm {
	if sStuGOrm == nil {
		logger.Error.Fatalln("must call CreateStuGOrm first")
		return nil
	}
	return sStuGOrm
}

func CreateStuGOrm() {
	var dbConfig DBConfig
	dbTree := config.GetConfig().GetRootTree().Get("database").(*toml.Tree)
	err := dbTree.Unmarshal(&dbConfig)
	if err != nil {
		logger.Error.Fatalln(err)
		return
	}
	var dbType = dbConfig.Type
	var creator = GetCreatorByType(dbType)
	if creator == nil {
		logger.Error.Fatalf("fail to find creator for type:%s", dbType)
		return
	}
	db, err := creator.Create(dbConfig)
	if err != nil {
		logger.Error.Fatalln(err)
		return
	}
	sStuGOrm = &StuGOrm{mDB: db}
}

func (sgo StuGOrm) GetDB() *gorm.DB {
	return sgo.mDB
}
