-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS user_response;

CREATE TABLE submission (
    id SERIAL PRIMARY KEY,
    question_set_id INT REFERENCES question_set(id),
    user_id INT REFERENCES users(id),
    correct_count INT NOT NULL,
    answers_count INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE response (
    id SERIAL PRIMARY KEY,
    submission_id INT REFERENCES submission(id),
    question_id INT REFERENCES question(id),
    answer_id INT REFERENCES answer(id),
    is_correct BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE user_response (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    question_id INT NOT NULL REFERENCES question(id),
    question_set_id INT NOT NULL REFERENCES question_set(id),
    answer_id INT NOT NULL REFERENCES answer(id),
    is_correct BOOLEAN NOT NULL,
    answered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS response;
DROP TABLE IF EXISTS submission;
-- +goose StatementEnd