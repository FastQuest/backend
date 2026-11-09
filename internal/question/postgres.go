package question

import (
	"strconv"

	"gorm.io/gorm"
)

// applyFilters adds to the query the Joins and Where clauses of every filter
// requested in the DTO. The SQL is PostgreSQL specific (ILIKE, EXTRACT).
func applyFilters(query *gorm.DB, filter QuestionFilterDTO) *gorm.DB {
	if filter.Statement != "" {
		query = query.Where("statement ILIKE ?", "%"+filter.Statement+"%")
	}

	if filter.QuestionOption != "" {
		query = query.Where("id IN (SELECT question_id FROM question_option WHERE text ILIKE ?)", "%"+filter.QuestionOption+"%")
	}

	if filter.Subject != "" {
		query = query.
			Joins("INNER JOIN subject ON subject.id = question.subject_id").
			Where("subject.id = ?", filter.Subject)
	}

	if filter.Topic != "" {
		query = query.Where("topic = ?", filter.Topic)
	}

	if filter.Source != "" {
		query = query.Joins("JOIN source_exam_instance ON source_exam_instance.id = question.source_exam_instance_id").
			Where("source_exam_instance.source_id = ?", filter.Source)
	}

	if filter.SourceExamInstanceID != "" {
		query = query.Where("question.source_exam_instance_id = ?", filter.SourceExamInstanceID)
	}

	if filter.Year != "" {
		if yearInt, err := strconv.Atoi(filter.Year); err == nil {
			query = query.
				Joins("LEFT JOIN source_exam_instance ON source_exam_instance.id = question.source_exam_instance_id").
				Where("source_exam_instance.year = ? OR EXTRACT(YEAR FROM question.created_at) = ?",
					yearInt,
					yearInt,
				)
		}
	}

	if filter.List != "" {
		if isList, err := strconv.ParseBool(filter.List); err == nil {
			if isList {
				query = query.Where("question_set_id IS NOT NULL")
			} else {
				query = query.Where("question_set_id IS NULL")
			}
		}
	}

	return query
}

// applyIncludes preloads the relations requested by the include query param.
// Unknown relations are ignored.
func applyIncludes(includes []string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, include := range includes {
			switch include {
			case "source":
				db = db.Preload("SourceExamInstance.Source")
			case "question-options":
				db = db.Preload("QuestionOptions")
			case "user":
				db = db.Preload("User")
			case "subject":
				db = db.Preload("Subject")
			}
		}
		return db
	}
}
