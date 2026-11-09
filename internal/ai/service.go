package ai

import (
	"context"
	"encoding/json"
	"flashquest/internal/question"
	"flashquest/internal/questionoption"
	"flashquest/internal/questionset"
	"flashquest/pkg/models"
	"flashquest/pkg/sliceutil"
	"fmt"
	"log"
	"time"

	"google.golang.org/genai"
)

var aiClient *genai.Client

type Service struct {
	questionService       question.Service
	questionOptionService questionoption.Service
	questionSetRepository *questionset.Repository
}

func NewService(questionService question.Service, questionOptionService questionoption.Service, questionSetRepository *questionset.Repository) *Service {
	return &Service{
		questionService:       questionService,
		questionOptionService: questionOptionService,
		questionSetRepository: questionSetRepository,
	}
}

func InitGemini() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Println("chatGemini error:", err)
		log.Fatal(err)
	}

	aiClient = client
}

func chatGemini(message string) (string, error) {
	ctx := context.Background()

	result, err := aiClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(message),
		nil,
	)
	if err != nil {
		log.Println("chatGemini error:", err)
		return "", err
	}

	return result.Text(), nil
}

func genQuestion(text string) AIQuestionResponse {
	prompt := fmt.Sprintf(`
Você é um assistente que fala somente em JSON focado em criar questões com 4 alternativas. Não escreva texto normal. Não use markdown. Sempre responda no formato JSON:

{
statement: "Texto para questão",
answers: [
{
text: "Alternativa 1",
is_correct: false
},
{
text: "Alternativa 2",
is_correct: false
},
{
text: "Alternativa 3",
is_correct: true
},
{
text: "Alternativa 4",
is_correct: false
}
]
}

Seguindo o formato a cima cria uma questão sobre: %s
`, text)

	var q AIQuestionResponse
	aiResponse, _ := chatGemini(prompt)

	err := json.Unmarshal([]byte(aiResponse), &q)
	if err != nil {
		log.Println("Error on convert response")
	} else {
		log.Println("Successful Generated")
	}

	return q
}

func genQuestionSet(text string) AIQuestionSetResponse {
	prompt := fmt.Sprintf(`
Você é um assistente que fala somente em JSON focado em criar uma lista de questões com 10 questões com 4 alternativas cada. Não escreva texto normal. Não use markdown. Sempre responda no formato JSON:

{
"name": "Nome da lista",
"description": "Descrição da lista",
questions: [
{
statement: "Texto para questão",
answers: [
{
text: "Alternativa 1",
is_correct: false
},
{
text: "Alternativa 2",
is_correct: false
},
{
text: "Alternativa 3",
is_correct: true
},
{
text: "Alternativa 4",
is_correct: false
}
]
}
]
}

Seguindo o formato a cima cria uma lista com 10 questões sobre: %s
`, text)

	var questionSet AIQuestionSetResponse
	aiResponse, _ := chatGemini(prompt)

	err := json.Unmarshal([]byte(aiResponse), &questionSet)
	if err != nil {
		log.Println("Error on convert response")
	} else {
		log.Println("Successful Generated")
	}

	return questionSet
}

func formatQuestions(aiQuestions ...AIQuestionResponse) []question.CreateRequest {
	questions := make([]question.CreateRequest, 0, len(aiQuestions))
	for _, q := range aiQuestions {
		questions = append(questions, question.CreateRequest{
			Statement: q.Statement,
			SubjectID: 7,
			UserID:    2,
		})
	}
	return questions
}

func formatQuestionOption(aiAnswers ...AIAnswerResponse) []questionoption.CreateRequest {
	questionoptions := make([]questionoption.CreateRequest, 0, len(aiAnswers))
	for _, a := range aiAnswers {
		questionoptions = append(questionoptions, questionoption.CreateRequest{
			Text:      a.Text,
			IsCorrect: a.IsCorrect,
		})
	}
	return questionoptions
}

func (s *Service) addAIQuestion(aiQuestion AIQuestionResponse) error {
	created, err := s.questionService.SendQuestions(formatQuestions(aiQuestion))
	if err != nil {
		return fmt.Errorf("failed to create generated question: %w", err)
	}

	if len(created) == 0 {
		return fmt.Errorf("question service returned no created question")
	}

	_, err = s.questionOptionService.CreateQuestionOptions(
		fmt.Sprint(created[0].ID),
		formatQuestionOption(aiQuestion.Answers...),
	)
	if err != nil {
		return fmt.Errorf("failed to create generated question options: %w", err)
	}
	log.Println("Successful Question and QuestionOption Insert")
	return nil
}

func (s *Service) addAIQuestionSet(aiQuestionSet AIQuestionSetResponse) error {
	questionSet := models.QuestionSet{
		Name:        aiQuestionSet.Name,
		Description: aiQuestionSet.Description,
		UserID:      2,
		CreatedAt:   time.Now(),
		IsPrivate:   false,
		Type:        "list",
	}

	errSendQS := s.questionSetRepository.SendQuestionSets(&questionSet)
	if errSendQS != nil {
		return errSendQS
	}

	questions, errSendQ := s.questionService.SendQuestions(formatQuestions(aiQuestionSet.Questions...))
	if errSendQ != nil {
		return errSendQ
	}

	for i, q := range aiQuestionSet.Questions {
		_, err := s.questionOptionService.CreateQuestionOptions(
			fmt.Sprint(questions[i].ID),
			formatQuestionOption(q.Answers...),
		)
		if err != nil {
			return fmt.Errorf("failed to create generated question options: %w", err)
		}
	}

	questionSetQuestion := make([]models.QuestionSetQuestion, 0, len(questions))
	for i, q := range questions {
		questionSetQuestion = append(questionSetQuestion, models.QuestionSetQuestion{
			QuestionSetID: questionSet.ID,
			QuestionID:    int(q.ID),
			Position:      i + 1,
		})
	}

	errSendQSQ := s.questionSetRepository.SendQuestionSetQuestionInternal(sliceutil.PtrSlice(questionSetQuestion)...)
	if errSendQSQ != nil {
		return errSendQSQ
	}

	return nil
}
