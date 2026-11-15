package source

type SourceExamBody struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Edition uint   `json:"edition"`
	Phase   uint   `json:"phase"`
	Year    uint   `json:"year"`
}
