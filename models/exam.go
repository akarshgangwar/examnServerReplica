package models

import (
	"time"

	"gorm.io/gorm"
)

/** Exam Model with following fields:
ID - int
Name - string
Description - string
CreatedAt - time.Time
UpdatedAt - time.Time
CreatedBy - string
UpdatedBy - string
**/

type Exam struct {
	gorm.Model
	ID          int        `gorm:"primaryKey"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at,string,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at_at,string,omitempty"`
	CreatedBy   string     `json:"created_by"`
	UpdatedBy   string     `json:"updated_by"`
}

func (e *Exam) TableName() string {
	return "exam"
}

/** Question Model with following fields:
ID - int
Question Statement - string
Type - string
CreatedAt - time.Time
UpdatedAt - time.Time
CreatedBy - string
UpdatedBy - string
**/

type Question struct {
	gorm.Model
	ID                int        `gorm:"primaryKey"`
	ExamID            int        `gorm:"index"` // Add the index tag to improve query performance
	Exam              Exam       `gorm:"foreignKey:ExamID"`
	QuestionStatement string     `json:"question_statement"`
	Type              string     `json:"type"`
	CreatedAt         *time.Time `json:"created_at,string,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at_at,string,omitempty"`
	CreatedBy         string     `json:"created_by"`
	UpdatedBy         string     `json:"updated_by"`
}

func (q *Question) TableName() string {
	return "question"
}

/** Options Model with following fields:
ID - int
QuestionID - Foreign Key
text - string
is_correct - bool
CreatedAt - time.Time
UpdatedAt - time.Time
CreatedBy - string
UpdatedBy - string
**/

type Option struct {
	gorm.Model
	ID         int        `gorm:"primaryKey"`
	QuestionID int        `gorm:"index"`                 // Add the index tag to improve query performance
	Question   Question   `gorm:"foreignKey:QuestionID"` // Correct foreign key definition
	Text       string     `json:"value"` 	
	Type	   string	   `json:"type"`	 	
	IsCorrect  bool       `json:"is_correct"`
	CreatedAt  *time.Time `json:"created_at,string,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,string,omitempty"`
	CreatedBy  string     `json:"created_by"`
	UpdatedBy  string     `json:"updated_by"`
}

func (o *Option) TableName() string {
	return "option"
}


type Attempt struct {
	ID         int        `gorm:"primaryKey"` 
	UserID     string     `gorm:"index"`
	User	   User	      `gorm:"foreignKey:UserID"` 
	ExamID     int        `gorm:"index"` 
	Exam       Exam       `gorm:"foreignKey:ExamID"`
	StartTime  time.Time  `json:"start_time"` 
	EndTime    *time.Time `json:"end_time"` 
	CreatedAt  *time.Time `json:"created_at,string,omitempty"` 
	UpdatedAt  *time.Time `json:"updated_at,string,omitempty"` 
}
type AttemptSerializer struct {
	ExamID     int        `json:"Exam_id"`
	ID			int		  `json:"Attempt_id"`
	StartTime  time.Time  `json:"start_time"` 
}
type AttemptLog struct {
	ID          int        `gorm:"primaryKey"` 
	AttemptID   int        `gorm:"index"` // Foreign key to Attempt table 
	Attempt 	Attempt	   `gorm:"foreignKey:AttemptID"` 
	QuestionID  int        `gorm:"index"` // Foreign key to Question table 
	Question   Question    `gorm:"foreignKey:QuestionID"` 
	Answer      string     `json:"answer"` 
	AnsweredAt  time.Time  `json:"answered_at"` 
	CreatedAt   *time.Time `json:"created_at,string,omitempty"` 
	UpdatedAt   *time.Time `json:"updated_at,string,omitempty"` 
}

type UserEvent struct {
	ID        int        `gorm:"primaryKey"`
	AttemptID int        `gorm:"index"` 
	Attempt 	Attempt	   `gorm:"foreignKey:AttemptID"`
	Timestamp time.Time  `json:"timestamp"`
	Payload   string     `json:"payload"`
	Type      string     `json:"type"`
	CreatedAt *time.Time `json:"created_at,string,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,string,omitempty"`
}

type QuestionFormat struct {
	ID       string   `json:"id"`
	Question QuestionSerializer `json:"question"`
	Options  []OptionSerializer `json:"options"`
}

type OptionSerializer  struct {
	// ID        int    `json:"id"`
	Type	  string `json:"type"`
	Value     string `json:"value"`
}

type QuestionSerializer  struct {
	ID        int        `gorm:"primaryKey"`
	Statement string     `json:"statement"`
	Type      string     `json:"type"`
	Options []OptionSerializer `json:"options"`
}

type SavedAnswer struct {
	ID         int        `gorm:"primaryKey"`
	AttemptID  int        `gorm:"index"` 
	Attempt    Attempt    `gorm:"foreignKey:AttemptID"`
	SavedAnswers string   `json:"saved_answers"`
}
