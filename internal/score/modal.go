package score

import "time"

type Score struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string    `gorm:"size:100;not null;column:firstname" json:"firstname"`
	LastName  string    `gorm:"size:100;not null;column:lastname" json:"lastname"`
	Email     string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Kind      string    `gorm:"not null" json:"kind"`
	Score     string    `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ScoreStruct struct {
	Score     string `json:"score"`
	Kind      string `json:"kind"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

func (Score) TableName() string {
	return "scores"
}
