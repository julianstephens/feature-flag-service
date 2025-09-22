package seeder

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/julianstephens/feature-flag-service/internal/rbac"
)

type Seeder struct {
	conn *pgx.Conn
}

func NewSeeder(conn *pgx.Conn) *Seeder {
	return &Seeder{conn: conn}
}

// SeedRoles inserts roles if not present
func (s *Seeder) SeedRoles(ctx context.Context, roles []rbac.CreateRbacRoleRequest) error {
	for _, role := range roles {
		var exists bool
		err := s.conn.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM rbac_roles WHERE name=$1)`, role.Name).
			Scan(&exists)
		if err != nil {
			return fmt.Errorf("check role %s existence: %w", role.Name, err)
		}
		if !exists {
			_, err := s.conn.Exec(ctx,
				`INSERT INTO rbac_roles (name, description, created_at, updated_at)
				 VALUES ($1, $2, $3, $3)`,
				role.Name, role.Description, time.Now(),
			)
			if err != nil {
				return fmt.Errorf("insert role %s: %w", role.Name, err)
			}
		}
	}
	return nil
}

// SeedPermissions inserts permissions if not present
func (s *Seeder) SeedPermissions(ctx context.Context, permissions []rbac.CreateRbacPermissionRequest) error {
	for _, perm := range permissions {
		var exists bool
		err := s.conn.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM rbac_permissions WHERE name=$1)`, perm.Name).
			Scan(&exists)
		if err != nil {
			return fmt.Errorf("check permission %s existence: %w", perm.Name, err)
		}
		if !exists {
			_, err := s.conn.Exec(ctx,
				`INSERT INTO rbac_permissions (name, description)
				 VALUES ($1, $2)`,
				perm.Name, perm.Description,
			)
			if err != nil {
				return fmt.Errorf("insert permission %s: %w", perm.Name, err)
			}
		}
	}
	return nil
}

// SeedRolePermissions maps roles to permissions
func (s *Seeder) SeedRolePermissions(ctx context.Context, mappings []rbac.RbacRolePermission) error {
	for _, mapping := range mappings {
		var roleID, permID int
		err := s.conn.QueryRow(ctx, `SELECT id FROM rbac_roles WHERE name=$1`, mapping.RoleName).
			Scan(&roleID)
		if err != nil {
			return fmt.Errorf("lookup role '%s': %w", mapping.RoleName, err)
		}
		err = s.conn.QueryRow(ctx, `SELECT id FROM rbac_permissions WHERE name=$1`, mapping.PermissionName).
			Scan(&permID)
		if err != nil {
			return fmt.Errorf("lookup permission '%s': %w", mapping.PermissionName, err)
		}
		var exists bool
		err = s.conn.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM rbac_role_permissions WHERE role_id=$1 AND permission_id=$2)`,
			roleID, permID).
			Scan(&exists)
		if err != nil {
			return fmt.Errorf("check mapping %s->%s existence: %w", mapping.RoleName, mapping.PermissionName, err)
		}
		if !exists {
			_, err = s.conn.Exec(ctx,
				`INSERT INTO rbac_role_permissions (role_id, permission_id) VALUES ($1, $2)`,
				roleID, permID)
			if err != nil {
				return fmt.Errorf("insert role-permission mapping %s->%s: %w", mapping.RoleName, mapping.PermissionName, err)
			}
		}
	}
	return nil
}
