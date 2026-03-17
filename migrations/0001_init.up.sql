CREATE TABLE IF NOT EXISTS monitors (
    id SERIAL PRIMARY KEY,
    url VARCHAR(100) NOT NULL,
    interval INT NOT NULL,
    status VARCHAR(30) DEFAULT 'unknown',
    last_check TIMESTAMPTZ NULL,
    next_check TIMESTAMPTZ NULL,
    response_time BIGINT NULL
);

CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    monitor_id INT NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    status VARCHAR(30) NOT NULL,
    check_time TIMESTAMPTZ NOT NULL,
    response_time BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_notes_monitor ON notes (monitor_id);