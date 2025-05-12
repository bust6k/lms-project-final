package application

import (
	"github.com/gin-gonic/gin"
	"project_yandex_lms/lms-calc-with-gorutine/config"

	"project_yandex_lms/lms-calc-with-gorutine/pkg/Orchestrator/web"
	"project_yandex_lms/lms-calc-with-gorutine/pkg/Orchestrator/web/auth"
)

type Application struct {
	cfg    *config.Config
	router *gin.Engine
}

func New() *Application {
	return &Application{
		cfg:    config.DefaultConfig(),
		router: gin.Default(),
	}
}

func (a *Application) Setup() {

	a.router.POST("/internal", web.ProcessTasksAndSubmitTasks)

	a.router.Any("/internal/task", web.GetOrPostTask)

	apiV1Group := a.router.Group("/api/v1")
	{

		apiV1Group.POST("/calculate", auth.AutentificationMiddleware(), auth.AutorizationMiddleWare(), web.CreateAUserExpressionInSystem)
		apiV1Group.Any("/expressions", auth.AutentificationMiddleware(), auth.AutorizationMiddleWare(), web.GetUserExpressions)
		apiV1Group.GET("/expressions/:id", auth.AutentificationMiddleware(), auth.AutorizationMiddleWare(), web.GetASeparateExpression)
		apiV1Group.Any("/register", auth.RegistrateUser)
		apiV1Group.Any("login", auth.AutorizateUser)

	}
}

func (a *Application) Run() {

	a.router.Run(a.cfg.Port)
}
