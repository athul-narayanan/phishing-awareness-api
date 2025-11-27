package score

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScoreController struct {
	ScoreService *ScoreService
}

func (sc *ScoreController) CreateScore(c *gin.Context) {
	var req struct {
		FirstName string `json:"firstname" binding:"required"`
		LastName  string `json:"lastname" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Score     string `json:"score" binding:"required"`
		Kind      string `json:"kind" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	score := Score{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Kind:      req.Kind,
		Score:     req.Score,
	}

	newScore, err := sc.ScoreService.CreateScore(score)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Score added successfully, Please visit the score tab to view your scores.",
		"user": map[string]interface{}{
			"id":        newScore.ID,
			"firstname": newScore.FirstName,
			"lastname":  newScore.LastName,
			"email":     newScore.Email,
			"kind":      newScore.Kind,
		},
	})
}

func (sc *ScoreController) GetScoresByUser(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email query param required"})
		return
	}

	scores, err := sc.ScoreService.GetScoresByUser(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert DB model to AI model
	var arr []ScoreStruct
	for _, s := range scores {
		arr = append(arr, ScoreStruct{
			Score:     s.Score,
			Kind:      s.Kind,
			FirstName: s.FirstName,
			LastName:  s.LastName,
		})
	}

	rec, err := sc.ScoreService.GetPhishingHistoryRecommendation(arr)
	if err != nil {
		rec = "Keep practicing phishing awareness to improve!"
	}

	c.JSON(http.StatusOK, gin.H{
		"scores":         scores,
		"recommendation": rec,
	})
}
