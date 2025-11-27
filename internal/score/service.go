package score

import (
	"context"
	"fmt"
	"phishing-awareness-api/config"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type ScoreService struct {
	DB  *gorm.DB
	CFG *config.Config
}

func (s *ScoreService) CreateScore(score Score) (*Score, error) {

	if err := s.DB.Create(&score).Error; err != nil {
		return nil, err
	}

	return &score, nil
}

func (s *ScoreService) GetScoresByUser(email string) ([]Score, error) {
	var scores []Score
	result := s.DB.Where("email = ?", email).Find(&scores)
	if result.Error != nil {
		return nil, result.Error
	}
	return scores, nil
}

func (s *ScoreService) GetPhishingHistoryRecommendation(scores []ScoreStruct) (string, error) {
	if len(scores) == 0 {
		return "Play a quiz or game to get personalized phishing safety tips.", nil
	}

	var scoreSummary []string
	name := fmt.Sprintf("%s %s", scores[0].FirstName, scores[0].LastName)

	for _, s := range scores {
		scoreSummary = append(scoreSummary, fmt.Sprintf("%s: %s", strings.Title(s.Kind), s.Score))
	}
	scoreInfo := strings.Join(scoreSummary, ", ")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(s.CFG.GeminiAPIKey))
	if err != nil {
		return "", fmt.Errorf("Gemini init error: %v", err)
	}

	model := client.GenerativeModel("gemini-2.0-flash")

	prompt := fmt.Sprintf(`
You are a cybersecurity expert. The user (%s) completed Several phishing awareness quiz and games.

Here are his/her score over time.
%s

Analyze:
- Trend (Improving or declining)
- Consistency within scores
Finally respond with:
- 1 Sentence summarizing their progress
- 2-3 actions that can be taken in bullet points
- Suggest some links to refer for week users
- Do not use positive words if inconsistency or decline found
- 6-7 sentences maximum including recommendation links
- Say like you are giving advice to same person
`, name, scoreInfo)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("Gemini API error: %v", err)
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}
