ALTER TABLE rbac_roles 
DROP COLUMN created_at,
DROP COLUMN updated_at;

ALTER TABLE rbac_roles
DROP COLUMN public_role_id;