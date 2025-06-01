package share

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"xiaohaiyun/internal/models/share"
	"xiaohaiyun/internal/utils/shareUrlUtils"
)

func UrlStatus(c *gin.Context) {
	var req share.UrlShareString
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Invalid JSON: " + err.Error(),
		})
		return
	}

	// 解析 UUID
	parsedUUID, err := uuid.Parse(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Invalid UUID: " + err.Error(),
		})
		return
	}

	// 构建 UrlShare 对象
	urlShare := share.UrlShare{
		ID:        parsedUUID[:],
		Url:       req.Url,
		Username:  req.Username,
		Signature: req.Signature,
		Email:     req.Email,
		UserReqID: req.UserReqID,
		Avatar:    req.Avatar,
	}

	// 业务逻辑
	status, err := shareUrlUtils.CompareData(&urlShare)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Internal server error 服务器错误",
		})
		return
	}

	if status {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "success",
			"data": urlShare,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "数据错误",
		})
	}
}
