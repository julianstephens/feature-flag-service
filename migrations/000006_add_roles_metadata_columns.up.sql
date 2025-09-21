-- Add created_at and updated_at columns with default timestamps
ALTER TABLE rbac_roles
ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- Add public_role_id column
ALTER TABLE rbac_roles
ADD COLUMN public_role_id UUID;