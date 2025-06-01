package file

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"strconv"
	cosFile "xiaohaiyun/internal/utils/cos"
)

type SrcType struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

// HTML 要用数据万象
func HTML(c *gin.Context) {
	userID := cosFile.GetID(c).ID
	client := cosFile.Client()
	var Src SrcType

	c.ShouldBind(&Src)

	Path := "users" + "/" + strconv.Itoa(userID) + "/" + Src.Key
	opt := &cos.DocPreviewHTMLOptions{
		DstType:  "html",
		SrcType:  Src.Type,
		Copyable: "1",
		HtmlParams: &cos.HtmlParams{
			CommonOptions: &cos.HtmlCommonParams{
				IsShowTopArea: false,
			},
			PptOptions: &cos.HtmlPptParams{
				IsShowBottomStatusBar: true,
			},
		},
		Htmlwaterword:  "5pWw5o2u5LiH6LGhLeaWh+aho+mihOiniA==",
		Htmlfillstyle:  "cmdiYSgxMDIsMjA0LDI1NSwwLjMp", // rgba(102,204,255,0.3)
		Htmlfront:      "Ym9sZCAyNXB4IFNlcmlm",         // bold 25px Serif
		Htmlrotate:     "315",
		Htmlhorizontal: "50",
		Htmlvertical:   "100",
	}
	resp, err := client.CI.DocPreviewHTML(context.Background(), Path, opt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": resp,
	})
}
