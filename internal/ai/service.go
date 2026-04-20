package ai

import (
	"context"
	"encoding/json"
	"flashquest/helpers"
	"flashquest/internal/answer"
	"flashquest/internal/question"
	"flashquest/internal/questionset"
	"flashquest/pkg/models"
	"fmt"
	"log"
	"time"

	"google.golang.org/genai"
)

var aiClient *genai.Client

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

func formatQuestions(aiQuestions ...AIQuestionResponse) []models.Question {
	questions := make([]models.Question, 0, len(aiQuestions))
	for _, q := range aiQuestions {
		questions = append(questions, models.Question{
			Statement: q.Statement,
			SubjectID: 7,
			UserID:    5,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
	return questions
}

func formatAnswer(questionID uint, aiAnswers ...AIAnswerResponse) []models.Answer {
	answers := make([]models.Answer, 0, len(aiAnswers))
	for _, a := range aiAnswers {
		answers = append(answers, models.Answer{
			Text:       a.Text,
			Is_correct: a.IsCorrect,
			QuestionID: questionID,
		})
	}
	return answers
}

func addAIQuestion(aiQuestion AIQuestionResponse) {
	q := formatQuestions(aiQuestion)[0]
	question.SendQuestions(&q)
	log.Println("Successful Question Insert")

	answersPayload := formatAnswer(q.ID, aiQuestion.Answers...)
	answer.SendAnswers(&answersPayload)
	log.Println("Successful Answer Insert")
}

func addAIQuestionSet(aiQuestionSet AIQuestionSetResponse) error {
	questionSet := models.QuestionSet{
		Name:        aiQuestionSet.Name,
		Description: aiQuestionSet.Description,
		UserID:      5,
		CreatedAt:   time.Now(),
		IsPrivate:   false,
		Type:        "list",
	}

	errSendQS := questionset.SendQuestionSets(&questionSet)
	if errSendQS != nil {
		return errSendQS
	}

	questions := formatQuestions(aiQuestionSet.Questions...)
	errSendQ := question.SendQuestions(helpers.PtrSlice(questions)...)
	if errSendQ != nil {
		return errSendQ
	}

	answers := make([]models.Answer, 0, len(aiQuestionSet.Questions)*4)
	for i, q := range aiQuestionSet.Questions {
		answers = append(answers, formatAnswer(questions[i].ID, q.Answers...)...)
	}

	errSendA := answer.SendAnswers(&answers)
	if errSendA != nil {
		return errSendA
	}

	questionSetQuestion := make([]models.QuestionSetQuestion, 0, len(questions))
	for i, q := range questions {
		questionSetQuestion = append(questionSetQuestion, models.QuestionSetQuestion{
			QuestionSetID: questionSet.ID,
			QuestionID:    int(q.ID),
			Position:      i + 1,
		})
	}

	errSendQSQ := questionset.SendQuestionSetQuestionInternal(helpers.PtrSlice(questionSetQuestion)...)
	if errSendQSQ != nil {
		return errSendQSQ
	}

	return nil
}
