-- +goose Up
-- +goose StatementBegin

ALTER TABLE answer RENAME TO question_option;
ALTER TABLE response RENAME COLUMN answer_id TO question_option_id;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE question_option RENAME TO answer;
ALTER TABLE response RENAME COLUMN question_option_id TO answer_id;

-- +goose StatementEnd