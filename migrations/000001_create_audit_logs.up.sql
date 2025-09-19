CREATE TABLE audit_logs (
    id UUID PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    user_id TEXT NOT NULL,
    action TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id TEXT NOT NULL,
    before JSONB,
    after JSONB,
    description TEXT
);