package question

import (
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func dryRunDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: "host=localhost"}), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatalf("failed to open dry run db: %v", err)
	}
	return db
}

func filteredSQL(t *testing.T, filter QuestionFilterDTO) string {
	t.Helper()
	db := dryRunDB(t)
	var questions []Question
	stmt := applyFilters(db.Model(&Question{}), filter).Find(&questions).Statement
	return db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
}

func TestApplyFiltersWithoutFilters(t *testing.T) {
	got := filteredSQL(t, QuestionFilterDTO{})
	if got != `SELECT * FROM "question"` {
		t.Fatalf("unexpected sql %s", got)
	}
}

func TestApplyFiltersBuildsEachClause(t *testing.T) {
	cases := []struct {
		name   string
		filter QuestionFilterDTO
		want   []string
	}{
		{"statement", QuestionFilterDTO{Statement: "lei"}, []string{`statement ILIKE '%lei%'`}},
		{"question option", QuestionFilterDTO{QuestionOption: "penal"}, []string{`id IN (SELECT question_id FROM question_option WHERE text ILIKE '%penal%')`}},
		{"subject", QuestionFilterDTO{Subject: "3"}, []string{`INNER JOIN subject ON subject.id = question.subject_id`, `subject.id = '3'`}},
		{"topic", QuestionFilterDTO{Topic: "x"}, []string{`topic = 'x'`}},
		{"source", QuestionFilterDTO{Source: "1"}, []string{`JOIN source_exam_instance ON source_exam_instance.id = question.source_exam_instance_id`, `source_exam_instance.source_id = '1'`}},
		{"year", QuestionFilterDTO{Year: "2024"}, []string{`LEFT JOIN source_exam_instance ON source_exam_instance.id = question.source_exam_instance_id`, `source_exam_instance.year = 2024 OR EXTRACT(YEAR FROM question.created_at) = 2024`}},
		{"list true", QuestionFilterDTO{List: "true"}, []string{`question_set_id IS NOT NULL`}},
		{"list false", QuestionFilterDTO{List: "false"}, []string{`question_set_id IS NULL`}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := filteredSQL(t, tc.filter)
			for _, fragment := range tc.want {
				if !strings.Contains(got, fragment) {
					t.Fatalf("expected %q in %s", fragment, got)
				}
			}
		})
	}
}

func TestApplyFiltersIgnoresUnparsableValues(t *testing.T) {
	got := filteredSQL(t, QuestionFilterDTO{Year: "abc", List: "maybe"})
	if got != `SELECT * FROM "question"` {
		t.Fatalf("unexpected sql %s", got)
	}
}
