package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	model2 "stu-manager/controller/model"
	"sync"
)

type Controller interface {
	Method() string
	RelativePath() string
	OnRequest(c *gin.Context)
}

func AddController(engine *gin.Engine, controller Controller) {
	engine.Handle(controller.Method(), controller.RelativePath(), controller.OnRequest)
}

const (
	ERR_INVALID_PARAM = "invalid_param"
	ERR_BAD_JSON      = "bad_json"
	ERR_DB_ERR        = "db_err"
)

func BindJsonOrError(c *gin.Context, request interface{}) bool {
	err := c.BindJSON(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, model2.CreateResponseFromError(ERR_BAD_JSON, err))
		return false
	}
	return true
}

var sDBLock sync.RWMutex //全局controller数据库锁，防止幻读
func AtomicDBOperationOrError(c *gin.Context, op func() error) bool {
	sDBLock.Lock()
	defer sDBLock.Unlock()
	err := op()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model2.CreateResponseFromError(ERR_DB_ERR, err))
		return false
	}
	return true
}

func ParseIntFromParamOrError(c *gin.Context, paramKey string) (bool, int64) {
	var valueStr = c.Param(paramKey)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			model2.CreateResponseFromError(ERR_INVALID_PARAM, err))
		return false, value
	}
	return true, value
}
