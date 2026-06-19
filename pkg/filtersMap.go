package filters

import (
	"strconv"

	"gorm.io/gorm"
)

type FilterFunc func(string, *gorm.DB) *gorm.DB

var QuestionFilters = map[string]FilterFunc{
	"statement": func(value string, qb *gorm.DB) *gorm.DB {
		return qb.Where("statement ILIKE ?", "%"+value+"%")
	},
	"question_option": func(value string, qb *gorm.DB) *gorm.DB {
		return qb.Where("id IN (SELECT question_id FROM question_option WHERE text ILIKE ?)", "%"+value+"%")
	},
	"subject": func(value string, qb *gorm.DB) *gorm.DB {
		return qb.
			Joins("INNER JOIN subject ON subject.id = question.subject_id").
			Where("subject.id = ?", value)
	},
	"topic": func(value string, qb *gorm.DB) *gorm.DB {
		return qb.Where("topic = ?", value)
	},
	"source": func(value string, qb *gorm.DB) *gorm.DB {
		return qb.Joins("JOIN source_exam_instance ON source_exam_instance.id = question.source_exam_instance_id").
			Where("source_exam_instance.source_id = ?", value)
	},
	"year": func(value string, qb *gorm.DB) *gorm.DB {
		if yearInt, err := strconv.Atoi(value); err == nil {
			return qb.
				Joins("LEFT JOIN source_exam_instance ON source_exam_instance.id = question.source_exam_instance_id").
				Where("source_exam_instance.year = ? OR EXTRACT(YEAR FROM question.created_at) = ?",
					yearInt,
					yearInt,
				)
		}
		return qb
	},
	"list": func(value string, qb *gorm.DB) *gorm.DB {
		if isList, err := strconv.ParseBool(value); err == nil {
			if isList {
				return qb.Where("question_set_id IS NOT NULL")
			}
			return qb.Where("question_set_id IS NULL")
		}
		return qb
	},
}
