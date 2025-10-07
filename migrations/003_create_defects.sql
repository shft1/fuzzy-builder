CREATE TABLE IF NOT EXISTS defects (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    assigned_to BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,
    status TEXT NOT NULL CHECK (status IN ('new','in_progress','on_review','closed')),
    priority TEXT NOT NULL CHECK (priority IN ('low','medium','high')),
    due_date TIMESTAMPTZ NULL,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


