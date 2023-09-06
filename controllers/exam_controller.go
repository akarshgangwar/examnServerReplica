package controllers

import (
	"examn_go/helpers"
	"examn_go/models"
	"examn_go/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func StartTesHandler(ctx *gin.Context) {
	examID, err := strconv.Atoi(ctx.Param("exam_id")) //extracting params
	email := ctx.GetString("email")
	var attempts []models.Attempt
	// check if user already has an attempt or not
	err = repository.GetAttemptsByEmail(&attempts, email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attempts"})
		return
	}
	// check if provided exam id is valid or not
	var exams models.Exam
	_, err = repository.GetExam(&exams, examID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No such Exam with given Exam ID"})
		return
	}

	// if no attempt found then create one
	createAttempt := new(models.Attempt)
	if len(attempts) == 0 {
		createAttempt.ExamID = examID
		createAttempt.StartTime = time.Now()
		createAttempt.UserID = email
		repository.Save(&createAttempt)
	}
	// User has an existing attempt
	if len(attempts) > 0 {
		var attemptSerializers []models.AttemptSerializer
		for _, attempt := range attempts {
			attemptSerializer := models.AttemptSerializer{
				ExamID:    attempt.ExamID,
				StartTime: attempt.StartTime,
				ID:        attempt.ID,
			}
			attemptSerializers = append(attemptSerializers, attemptSerializer)
		}
		ctx.JSON(http.StatusForbidden, gin.H{"has_existing_attempt": true, "attempt_info": attemptSerializers})
		return
	}
	// User does not have an existing attempt
	res, err := helpers.GetQuestionsAndOptions(examID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": "Failed to fetch questions and options"})
	}
	ctx.JSON(http.StatusOK, gin.H{"has_existing_attempt ": "false", "attempt_id ": createAttempt.ID, "questions_options ": res})

}

func GetOngoingAttempt(ctx *gin.Context) {
	email := ctx.GetString("email")

	attemptID, err := repository.GetOngoingAttemptID(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for ongoing attempt"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"attempt_id": attemptID})
}

func ReloadAttempt(ctx *gin.Context) {
	attemptIDStr := ctx.Param("attempt_id")
	attemptID, err := strconv.Atoi(attemptIDStr)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid attempt ID"})
		return
	}
	attempt, err := repository.GetAttemptByID(attemptID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid attempt ID"})
		return
	}
	res, err := helpers.GetQuestionsAndOptions(attempt.ExamID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": "Failed to fetch questions and options"})
	}

	savedAnswers, markedForReview, err := repository.GetSavedAnswersAndMarkedForReview(attemptID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"start_time": attempt.StartTime, "questions": res, "saved_answers": savedAnswers, "marked_for_review": markedForReview})
}

type SaveAnswerRequest struct {
	// AttemptID  int        `json:"attempt_id`
	SavedAnswers []string `json:"saved_answers"`
}

func EndTest(ctx *gin.Context) {
	attemptIDStr := ctx.Param("attempt_id")
	attemptID, err := strconv.Atoi(attemptIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attempt ID"})
		return
	}
	var req SaveAnswerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	savedAnswer := make([]models.SavedAnswer, 0, len(req.SavedAnswers))
	for _, ans := range req.SavedAnswers {
		answer := models.SavedAnswer{
			AttemptID:    attemptID,
			SavedAnswers: ans,
		}
		savedAnswer = append(savedAnswer, answer)
	}
	err = repository.SaveAnswers(savedAnswer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save answers"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Answers successfully recorded"})
}
