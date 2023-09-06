package controllers

import (
	"examn_go/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllQuestions(ctx *gin.Context) {
	examID, err := strconv.Atoi(ctx.Param("exam_id")) //extracting params
	res, err := helpers.GetQuestionsAndOptions(examID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": "Failed to fetch questions and options"})
	}
	ctx.JSON(http.StatusOK, gin.H{"questions":res})
}
