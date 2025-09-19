CREATE TABLE rbac_user_roles (
    user_id UUID NOT NULL REFERENCES rbac_users(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES rbac_roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id)
);