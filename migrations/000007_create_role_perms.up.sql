CREATE TABLE rbac_permissions (
  id SERIAL PRIMARY KEY,
  name VARCHAR(64) UNIQUE NOT NULL,
  description TEXT
);

CREATE TABLE rbac_role_permissions (
  role_id INT REFERENCES rbac_roles(id) ON DELETE CASCADE,
  permission_id INT REFERENCES rbac_permissions(id) ON DELETE CASCADE,
  PRIMARY KEY (role_id, permission_id)
);
