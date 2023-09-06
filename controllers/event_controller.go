package controllers

import (
	"examn_go/models"
	"examn_go/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func PostUserEvents(ctx *gin.Context) {
	attemptIDStr := ctx.Param("attempt_id")
	attemptID, err := strconv.Atoi(attemptIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attempt ID"})
		return
	}
	var events []models.UserEvent
	if err := ctx.ShouldBindJSON(&events); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
		return
	}
	// Validate the attempt ID and check if the test is still ongoing (not ended)
	_, err = repository.GetAttemptByID(attemptID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid attempt ID or attempt has ended"})
		return
	}
	userEvents := make([]models.UserEvent, 0, len(events))
	timestamp := time.Now()
	for _, event := range events {
		userEvent := models.UserEvent{
			AttemptID: attemptID,
			Timestamp: timestamp,
			Payload:   event.Payload,
			Type:      event.Type,
		}
		userEvents = append(userEvents, userEvent)
	}
	// Save the events to the database at once
	err = repository.SaveUserEvents(userEvents)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save events"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Events successfully recorded"})
}
