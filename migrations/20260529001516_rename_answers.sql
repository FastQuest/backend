-- +goose Up
-- +goose StatementBegin

ALTER TABLE answer RENAME TO question_option;
ALTER TABLE submission RENAME COLUMN hits TO correct_count;
ALTER TABLE response RENAME COLUMN answer_id TO question_option_id;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE question_option RENAME TO answer;
ALTER TABLE submission RENAME COLUMN correct_count TO hits;
ALTER TABLE response RENAME COLUMN question_option_id TO answer_id;

-- +goose StatementEnd