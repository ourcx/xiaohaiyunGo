package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"xiaohaiyun/internal/api"
	"xiaohaiyun/internal/api/Aichat"
	"xiaohaiyun/internal/api/advice"
	"xiaohaiyun/internal/api/chat"
	"xiaohaiyun/internal/api/file"
	D3 "xiaohaiyun/internal/api/file/D3Data"
	"xiaohaiyun/internal/api/file/Describe"
	recyclebin "xiaohaiyun/internal/api/file/RecycleBin"
	"xiaohaiyun/internal/api/file/Tag"
	"xiaohaiyun/internal/api/file/share"
	"xiaohaiyun/internal/api/file/share/Manage"
	"xiaohaiyun/internal/api/userData"
	"xiaohaiyun/internal/api/userData/relationship"
	"xiaohaiyun/internal/api/userData/userFound"
	"xiaohaiyun/internal/controllers"
	userAuth "xiaohaiyun/internal/middleware/user"
	"xiaohaiyun/internal/services"
	cosFile "xiaohaiyun/internal/utils/cos"
	legislations "xiaohaiyun/internal/utils/legislation"
	"xiaohaiyun/internal/utils/reqEmailSend"
)

func SetupRoutes(r *gin.Engine, engine *xorm.Engine) {
	// 定义用户路由组
	user := r.Group("/user")
	user.Use(userAuth.GetExitJwt)
	{
		// 创建 UserService 实例
		UserService := services.NewUserService(engine)
		// 创建 UserController 实例
		UserController := controllers.NewUserController(UserService)

		user.GET("/", UserController.GetUsers)
		user.POST("/userReq", UserReq)
		user.POST("/userLongin", userLongin)
		user.POST("/parseJwt", api.JwtStatus)
		user.POST("/sendEmail", reqEmailSend.SendReqEmail)
		user.POST("/search", relationship.SearchUser)

	}

	pwd := r.Group("/pwd")
	pwd.Use(userAuth.LoggerMiddleware, userAuth.AuthMiddleware)
	{
		UserService := services.NewUserService(engine)
		UserController := controllers.NewUserController(UserService)
		pwd.GET("/", UserController.GetUsers)
		//密码更改
		pwd.PUT("/passwordChange", PasswordChange)
	}

	Jwt := r.Group("/Jwt")
	Jwt.Use(userAuth.LoggerMiddleware, userAuth.AuthMiddleware)
	{
		UserService := services.NewUserService(engine)
		UserController := controllers.NewUserController(UserService)
		Jwt.GET("/", UserController.GetUsers)
		Jwt.PUT("/exit", userAuth.SetExitJwt)
		Jwt.POST("/advice", advice.SendAdvice)
	}

	chats := r.Group("/chat")
	{
		UserService := services.NewUserService(engine)
		UserController := controllers.NewUserController(UserService)

		chats.GET("/", UserController.GetUsers)
		// WebSocket路由
		chats.POST("/groupHistory", chat.SetGroupHistory)
		chats.GET("/getGroupHistory", chat.GetGroupHistory)
		chats.GET("/ws", chat.HandleWebSocket)
		// 启动后台消息处理
		go chat.HandleMessages()
		go chat.HistoryMessage()
		//chats.GET("/double", chat.HandlePrivateChat)

	}
	profiles := r.Group("/profiles")
	profiles.Use(userAuth.LoggerMiddleware, userAuth.AuthMiddleware)
	{
		UserService := services.NewUserService(engine)
		UserController := controllers.NewUserController(UserService)
		profiles.GET("/", UserController.GetUsers)

		profiles.GET("/GetProfile", userFound.GetProfiles)
		profiles.POST("/PostProfile", userFound.PostProfiles)
		profiles.GET("/conversations", chat.GetConversations)
		profiles.POST("/setConversation", chat.SetConversation)
		profiles.POST("/relationApply", relationship.ApplyByEmail)
		profiles.POST("/upUserName", userFound.UpUserReqName)
		profiles.GET("/getRBook", relationship.GetRBookList)
		profiles.POST("/search", relationship.SearchFriend)
		//用户列表
		profiles.GET("/read", chat.Read)
		profiles.PUT("/forgetPwd", userData.ForgetPwd)

	}
	files := r.Group("/files")
	files.Use(userAuth.LoggerMiddleware, userAuth.AuthMiddleware)
	profiles.Use(userAuth.LoggerMiddleware, userAuth.AuthMiddleware)
	{
		UserService := services.NewUserService(engine)
		UserController := controllers.NewUserController(UserService)
		files.GET("/", UserController.GetUsers)
		files.GET("/init", file.Init)
		files.GET("/ListFile", cosFile.GenerateSecureUploadURL)
		files.GET("treeFIle", file.TreeFile)
		files.POST("/ReplayFile", file.ReplayList)
		files.POST("/removeFIleName", file.RemoveFile)
		files.POST("/removeFolder", file.RemoveFolder)
		files.POST("/Ourl", file.Ourl)
		files.POST("/Dscribe", Describe.Describe)
		files.POST("/ForDescribe", file.Delete)
		files.POST("/AddFolder", file.AddFolder)
		files.POST("/baseData", file.BaseData)
		files.PUT("/RenameFile", file.RenameFile)
		files.POST("/delete", recyclebin.Trash)
		files.POST("/move", file.MoveFile)
		files.POST("/copy", file.CopyFile)
		files.PUT("/special", file.SpecialTreeFile)
		files.POST("/imgData", file.ImgBaseData)
		files.GET("/imgDate", file.ImgDate)
		files.PUT("/setTag", Tag.SetTag)
		files.PUT("/getTag", Tag.GetTag)
		files.PUT("deleteTag", Tag.DeleteTag)
		files.POST("/HTML", file.HTML)
	}
	outShare := r.Group("/outShare")
	{
		UserService := services.NewUserService(engine)
		UserController := controllers.NewUserController(UserService)
		outShare.GET("/", UserController.GetUsers)
		outShare.POST("/urlStatus", share.UrlStatus)
		outShare.PUT("/access", share.AccessVisit)
		//传前端query的参数，可以进行校验，如果检验不通过，则通知前端退出当前状态

	}
	shareUrl := r.Group("/share")
	shareUrl.Use(userAuth.LoggerMiddleware, userAuth.AuthMiddleware)
	{
		UserService := services.NewUserService(engine)
		UserController := controllers.NewUserController(UserService)
		shareUrl.GET("/", UserController.GetUsers)
		shareUrl.POST("/setShareData", share.SetShare)
		shareUrl.POST("/getShareData", share.GetShare)
		shareUrl.POST("/create", share.CreateUrl)
		shareUrl.POST("/getUrl", share.GetUrl)
		shareUrl.POST("/checked", share.Checked)
		shareUrl.DELETE("/deleteShare", share.DeleteShare)
	}
	email := r.Group("/email")
	email.Use(userAuth.LoggerMiddleware, userAuth.AuthMiddleware)
	email.Use(legislations.StartCleanupScheduler)
	{
		UserService := services.NewUserService(engine)
		UserController := controllers.NewUserController(UserService)
		email.GET("/", UserController.GetUsers)
		email.GET("/sendEmail", legislations.SendEmail)
		email.POST("/requestReset", ValidateCode)
	}
	Data := r.Group("/data")
	Data.Use(userAuth.LoggerMiddleware, userAuth.AuthMiddleware)
	{
		UserService := services.NewUserService(engine)
		UserController := controllers.NewUserController(UserService)
		Data.GET("/", UserController.GetUsers)
		Data.GET("/total", D3.Total)
		Data.GET("/proportion", D3.Proportion)
		Data.GET("/shareData", Manage.ShareList)
		Data.GET("/logins", userFound.SendLog)
		Data.GET("/relationD3", D3.SearchRelationshipByID)
	}
	AI := r.Group("/AI")
	AI.Use(userAuth.LoggerMiddleware, userAuth.AuthMiddleware)
	{
		UserService := services.NewUserService(engine)
		UserController := controllers.NewUserController(UserService)
		AI.GET("/", UserController.GetUsers)
		AI.PUT("/chat", Aichat.AiChat)
	}
	trash := r.Group("/trash")
	trash.Use(userAuth.LoggerMiddleware, userAuth.AuthMiddleware)
	{
		UserService := services.NewUserService(engine)
		UserController := controllers.NewUserController(UserService)
		trash.GET("/", UserController.GetUsers)
		trash.POST("/addTrash", recyclebin.Trash)
		trash.GET("/TrashList", recyclebin.TrashList)
		trash.PUT("/RecoverFile", recyclebin.RecoverTrashFile)
		trash.PUT("/deleteTrash", recyclebin.DeleteTrashList)
	}
}
