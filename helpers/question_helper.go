package helpers

import (
	"examn_go/models"
	"examn_go/repository"

)

func GetQuestionsAndOptions(examId int) ([]models.QuestionSerializer, error) {
	var questions []models.Question
	var questionSerializers []models.QuestionSerializer
	err := repository.GetQuestionsByExamID(examId, &questions)
	if err != nil {
		return questionSerializers, err
	}
	var questionIDs []int
	for _, question := range questions {
		questionIDs = append(questionIDs, question.ID)
	}
	// Fetch options for each question
	var options []models.Option
	err = repository.GetOptionsByQuestionIDs(questionIDs, &options)
	if err != nil {
		return nil, err
	}
	// Create a map to associate options with their respective question IDs
	questionIDToOptions := make(map[int][]models.Option)
	for _, option := range options {
		questionIDToOptions[option.QuestionID] = append(questionIDToOptions[option.QuestionID], option)
	}
	// Create question serializers with associated options
	for _, question := range questions {
		optionSerializers := make([]models.OptionSerializer, len(questionIDToOptions[question.ID]))
		for i, option := range questionIDToOptions[question.ID] {
			optionSerializers[i] = models.OptionSerializer{
				Type:  option.Type,
				Value: option.Text,
			}
		}
		questionSerializers = append(questionSerializers, models.QuestionSerializer{
			ID:        question.ID,
			Statement: question.QuestionStatement,
			Type:      question.Type,
			Options:   optionSerializers,
		})
	}
	return questionSerializers, nil
}
