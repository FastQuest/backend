package classroom

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormRepository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &gormRepository{db: db} }

func (r *gormRepository) CreateClassroom(classroom *Classroom) error {
	return r.db.Create(classroom).Error
}

func (r *gormRepository) ListClassroomsForUser(userID uint) ([]Classroom, error) {
	var classrooms []Classroom
	err := r.db.Model(&Classroom{}).
		Joins("LEFT JOIN classroom_participant cp ON cp.classroom_id = classroom.id AND cp.user_id = ? AND cp.status = ?", userID, StatusActive).
		Where("classroom.creator_id = ? OR cp.id IS NOT NULL", userID).
		Order("classroom.created_at DESC").
		Find(&classrooms).Error
	return classrooms, err
}

func (r *gormRepository) UpdateClassroomStatus(classroomID uint, status string) error {
	result := r.db.Model(&Classroom{}).Where("id = ?", classroomID).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *gormRepository) AddParticipant(classroomID, userID uint) (*Participant, error) {
	participant := Participant{ClassroomID: classroomID, UserID: userID, Status: StatusActive}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var classroom Classroom
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&classroom, classroomID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if classroom.Status != StatusActive {
			return ErrInactiveClassroom
		}
		return tx.Create(&participant).Error
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrDuplicateParticipant
		}
		return nil, err
	}
	return &participant, nil
}

func (r *gormRepository) ListParticipants(classroomID uint) ([]Participant, error) {
	var participants []Participant
	err := r.db.Where("classroom_id = ?", classroomID).Order("created_at ASC").Find(&participants).Error
	return participants, err
}

func (r *gormRepository) CreateActivity(classroomID, questionSetID, createdBy uint, title string, allowPrivateQuestionSet bool) (*Activity, error) {
	var activity Activity
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var classroom Classroom
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&classroom, classroomID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if classroom.Status != StatusActive {
			return ErrInactiveClassroom
		}
		var questionSet struct {
			ID        uint
			UserID    uint
			IsPrivate bool
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Table("question_set").Select("id, user_id, is_private").First(&questionSet, questionSetID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInvalidQuestionSet
			}
			return err
		}
		if questionSet.IsPrivate && questionSet.UserID != createdBy && !allowPrivateQuestionSet {
			return ErrQuestionSetForbidden
		}

		type sourceQuestion struct {
			QuestionID uint
			Statement  string
		}
		var questions []sourceQuestion
		if err := tx.Table("question_set_question qsq").
			Select("q.id AS question_id, q.statement").
			Joins("JOIN question q ON q.id = qsq.question_id").
			Where("qsq.question_set_id = ?", questionSetID).
			Order("qsq.position ASC, q.id ASC").
			Scan(&questions).Error; err != nil {
			return err
		}
		if len(questions) == 0 {
			return ErrQuestionSetEmpty
		}

		var versionNumber int
		if err := tx.Model(&QuestionSetVersion{}).
			Where("question_set_id = ?", questionSetID).
			Select("COALESCE(MAX(version_number), 0)").
			Scan(&versionNumber).Error; err != nil {
			return err
		}
		version := QuestionSetVersion{QuestionSetID: questionSetID, VersionNumber: versionNumber + 1, CreatedBy: createdBy}
		if err := tx.Create(&version).Error; err != nil {
			return err
		}

		seenQuestionIDs := make(map[uint]struct{}, len(questions))
		for i, source := range questions {
			if _, duplicate := seenQuestionIDs[source.QuestionID]; duplicate {
				return ErrInvalidQuestionSet
			}
			seenQuestionIDs[source.QuestionID] = struct{}{}
			versionQuestion := VersionQuestion{
				VersionID: version.ID, SourceQuestionID: source.QuestionID,
				Position: i + 1, Statement: source.Statement,
			}
			if err := tx.Create(&versionQuestion).Error; err != nil {
				return err
			}
			type sourceOption struct {
				ID        uint
				Text      string
				IsCorrect bool
			}
			var options []sourceOption
			if err := tx.Table("question_option").
				Select("id, text, is_correct").
				Where("id_question = ?", source.QuestionID).
				Order("id ASC").Scan(&options).Error; err != nil {
				return err
			}
			if len(options) == 0 {
				return ErrInvalidQuestionSet
			}
			for i, option := range options {
				versionOption := VersionOption{
					VersionQuestionID: versionQuestion.ID, SourceOptionID: option.ID,
					Position: i + 1, Text: option.Text, IsCorrect: option.IsCorrect,
				}
				if err := tx.Create(&versionOption).Error; err != nil {
					return err
				}
			}
		}

		activity = Activity{
			ClassroomID: classroomID, VersionID: version.ID,
			Title: title, Status: StatusActive,
		}
		if err := tx.Create(&activity).Error; err != nil {
			return err
		}
		activity.Version.QuestionSetID = questionSetID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (r *gormRepository) GetSubmissionContext(activityID, userID uint) (*SubmissionContext, error) {
	var activity Activity
	err := r.db.Joins("JOIN classroom ON classroom.id = classroom_activity.classroom_id").
		Joins("JOIN classroom_participant cp ON cp.classroom_id = classroom.id AND cp.user_id = ? AND cp.status = ?", userID, StatusActive).
		Preload("Version.Questions.Options").
		Where("classroom_activity.id = ?", activityID).
		Where("classroom_activity.status = ? AND classroom.status = ?", StatusActive, StatusActive).
		First(&activity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotActiveParticipant
	}
	if err != nil {
		return nil, err
	}

	context := &SubmissionContext{
		ActivityID: activity.ID, QuestionSetID: activity.Version.QuestionSetID,
		Questions: make([]SnapshotQuestion, len(activity.Version.Questions)),
	}
	for i, question := range activity.Version.Questions {
		options := make(map[uint]bool, len(question.Options))
		for _, option := range question.Options {
			options[option.SourceOptionID] = option.IsCorrect
		}
		context.Questions[i] = SnapshotQuestion{ID: question.SourceQuestionID, Options: options}
	}
	return context, nil
}

func (r *gormRepository) ListActivitiesForUser(classroomID, userID uint, allowAdmin bool) ([]Activity, error) {
	var activities []Activity
	err := r.db.Model(&Activity{}).
		Joins("JOIN classroom ON classroom.id = classroom_activity.classroom_id").
		Joins("LEFT JOIN classroom_participant cp ON cp.classroom_id = classroom.id AND cp.user_id = ? AND cp.status = ?", userID, StatusActive).
		Where("classroom_activity.classroom_id = ?", classroomID).
		Where("classroom.status = ? AND (classroom.creator_id = ? OR cp.id IS NOT NULL OR ?)", StatusActive, userID, allowAdmin).
		Preload("Version").
		Preload("Version.Questions", func(db *gorm.DB) *gorm.DB { return db.Order("position ASC") }).
		Preload("Version.Questions.Options", func(db *gorm.DB) *gorm.DB { return db.Order("position ASC") }).
		Order("classroom_activity.created_at DESC").
		Find(&activities).Error
	return activities, err
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(message, "unique constraint") || strings.Contains(message, "duplicate key")
}
