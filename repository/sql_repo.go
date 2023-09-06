package repository

import (
	"examn_go/infra/database"
	"examn_go/infra/logger"
	"examn_go/models"
	"fmt"
)

func Save(model interface{}) interface{} {
	err := database.DB.Create(model).Error
	if err != nil {
		logger.Errorf("error, not save data %v", err)
	}
	return err
}

func Get(model interface{}) interface{} {
	err := database.DB.Find(model).Error
	return err
}

func GetOne(model interface{}) interface{} {
	err := database.DB.Last(model).Error
	return err
}

func Update(model interface{}) error {
	err := database.DB.Save(model).Error
	return err
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	database.DB = database.DB.Debug()
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func GetQuestionsByExamID(examID int, questions *[]models.Question) error {
	// Query the database to get all questions with the specified ExamID
	err := database.DB.Where("exam_id = ?", examID).Find(questions).Error
	if err != nil {
		return err
	}
	return nil
}

func GetOptionsByQuestionIDs(questionIDs []int, options *[]models.Option) error {
	// Query the database to get options for the specified question IDs
	err := database.DB.Where("question_id IN (?)", questionIDs).Find(options).Error
	if err != nil {
		return err
	}
	return nil
}
func GetAttemptsByEmail(attempts *[]models.Attempt, email string) error {
	err := database.DB.Debug().Where("user_id = ?", email).Find(attempts).Error
	if err != nil {
		return fmt.Errorf("failed to fetch attempts: %w", err)
	}
	return nil
}

func GetExam(model interface{}, examId int) (interface{}, error) {
	database.DB = database.DB.Debug()
	err := database.DB.Where("id = ?", examId).First(model).Error
	return model, err
}

// GetOngoingAttemptID retrieves the attempt ID for the ongoing test of a specific user.
func GetOngoingAttemptID(email string) (int, error) {
	var attemptID int
	err := database.DB.Model(&models.Attempt{}).Where("user_id = ? AND end_time IS NULL", email).Pluck("id", &attemptID).Error
	if err != nil {
		return 0, err
	}

	return attemptID, nil
}

func GetAttemptByID(attemptID int) (*models.Attempt, error) {
	var attempt models.Attempt
	err := database.DB.Where("id = ?", attemptID).First(&attempt).Error
	if err != nil {
		logger.Errorf("error, failed to fetch attempt with ID %d: %v", attemptID, err)
		return nil, err
	}
	return &attempt, nil
}

func SaveUserEvents(events []models.UserEvent) error {
	err := database.DB.Create(&events).Error
	if err != nil {
		return err
	}
	return nil
}

func GetSavedAnswersAndMarkedForReview(attemptID int) (savedAnswers []string, markedForReview []string, err error) {
	var payloads []string
	var events []models.UserEvent
	if err := database.DB.Model(&models.UserEvent{}).
		Where("attempt_id = ? AND (type = 'saved_answer' OR type = 'marked_for_review')", attemptID).
		Find(&events).
		Error; err != nil {
		return nil, nil, err
	}
	// Separate the events into savedAnswers and markedForReview slices
	fmt.Print(payloads)
	for _, event := range events {
		if event.Type == "saved_answer" {
			savedAnswers = append(savedAnswers, event.Payload)
		} else if event.Type == "marked_for_review" {
			markedForReview = append(markedForReview, event.Payload)
		}
	}
	return savedAnswers, markedForReview, nil
}

func SaveAnswers(answers []models.SavedAnswer) error {
	err := database.DB.Create(&answers).Error
	if err != nil {
		return err
	}
	return nil
}
