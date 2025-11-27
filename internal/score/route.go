package score

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, sc *ScoreService) {
	controller := &ScoreController{ScoreService: sc}

	userGroup := r.Group("/phishing_awareness_api/score")

	{
		userGroup.GET("", controller.GetScoresByUser)
		userGroup.POST("", controller.CreateScore)
	}

}
