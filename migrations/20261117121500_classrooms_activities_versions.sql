-- +goose Up
-- +goose StatementBegin
CREATE TABLE classroom (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    creator_id INT NOT NULL REFERENCES users(id),
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE classroom_participant (
    id SERIAL PRIMARY KEY,
    classroom_id INT NOT NULL REFERENCES classroom(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id),
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (classroom_id, user_id)
);

CREATE TABLE question_set_version (
    id SERIAL PRIMARY KEY,
    question_set_id INT NOT NULL REFERENCES question_set(id),
    version_number INT NOT NULL,
    created_by INT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (question_set_id, version_number)
);

CREATE TABLE question_set_version_question (
    id SERIAL PRIMARY KEY,
    version_id INT NOT NULL REFERENCES question_set_version(id) ON DELETE CASCADE,
    source_question_id INT NOT NULL,
    position INT NOT NULL,
    statement TEXT NOT NULL,
    UNIQUE (version_id, position),
    UNIQUE (version_id, source_question_id)
);

CREATE TABLE question_set_version_option (
    id SERIAL PRIMARY KEY,
    version_question_id INT NOT NULL REFERENCES question_set_version_question(id) ON DELETE CASCADE,
    source_option_id INT NOT NULL,
    position INT NOT NULL,
    text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL,
    UNIQUE (version_question_id, position),
    UNIQUE (version_question_id, source_option_id)
);

CREATE TABLE classroom_activity (
    id SERIAL PRIMARY KEY,
    classroom_id INT NOT NULL REFERENCES classroom(id) ON DELETE CASCADE,
    version_id INT NOT NULL REFERENCES question_set_version(id),
    title TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE submission
    ADD COLUMN activity_id INT REFERENCES classroom_activity(id);

CREATE INDEX idx_classroom_participant_user_status
    ON classroom_participant(user_id, status);
CREATE INDEX idx_classroom_activity_classroom_status
    ON classroom_activity(classroom_id, status);
CREATE INDEX idx_submission_activity_id ON submission(activity_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_submission_activity_id;
DROP INDEX IF EXISTS idx_classroom_activity_classroom_status;
DROP INDEX IF EXISTS idx_classroom_participant_user_status;
ALTER TABLE submission DROP COLUMN IF EXISTS activity_id;
DROP TABLE IF EXISTS classroom_activity;
DROP TABLE IF EXISTS question_set_version_option;
DROP TABLE IF EXISTS question_set_version_question;
DROP TABLE IF EXISTS question_set_version;
DROP TABLE IF EXISTS classroom_participant;
DROP TABLE IF EXISTS classroom;
-- +goose StatementEnd
