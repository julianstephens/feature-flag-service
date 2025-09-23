package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/julianstephens/feature-flag-service/internal/rbac"
	"github.com/julianstephens/feature-flag-service/internal/seeder"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		fmt.Fprintln(os.Stderr, "DB_URL env var is required")
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	s := seeder.NewSeeder(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Define your roles
	roles := []rbac.CreateRbacRoleRequest{
		{PublicRoleID: uuid.New().String(), Name: "admin", Description: "Full access to all resources and administrative functions."},
		{PublicRoleID: uuid.New().String(), Name: "editor", Description: "Can modify feature flags and configurations."},
		{PublicRoleID: uuid.New().String(), Name: "auditor", Description: "Read-only access to audit logs and system status."},
		{PublicRoleID: uuid.New().String(), Name: "user", Description: "Can read feature flags and config, but not modify."},
	}

	// Define your permissions
	permissions := []rbac.CreateRbacPermissionRequest{
		{Name: "flags.read", Description: "Read feature flags"},
		{Name: "flags.write", Description: "Create or update feature flags"},
		{Name: "flags.delete", Description: "Delete feature flags"},
		{Name: "config.read", Description: "Read configuration values"},
		{Name: "config.write", Description: "Update configuration values"},
		{Name: "audit.read", Description: "Read audit logs"},
		{Name: "rbac.read", Description: "Read RBAC users and roles"},
		{Name: "rbac.write", Description: "Modify RBAC users and roles"},
	}

	// Map roles to permissions
	mappings := []rbac.RbacRolePermission{
		// Admin gets all
		{RoleName: "admin", PermissionName: "flags.read"}, {RoleName: "admin", PermissionName: "flags.write"}, {RoleName: "admin", PermissionName: "flags.delete"},
		{RoleName: "admin", PermissionName: "config.read"}, {RoleName: "admin", PermissionName: "config.write"},
		{RoleName: "admin", PermissionName: "audit.read"}, {RoleName: "admin", PermissionName: "rbac.read"}, {RoleName: "admin", PermissionName: "rbac.write"},
		// Editor
		{RoleName: "editor", PermissionName: "flags.read"}, {RoleName: "editor", PermissionName: "flags.write"}, {RoleName: "editor", PermissionName: "config.read"}, {RoleName: "editor", PermissionName: "config.write"},
		// Auditor
		{RoleName: "auditor", PermissionName: "flags.read"}, {RoleName: "auditor", PermissionName: "audit.read"},
		// User
		{RoleName: "user", PermissionName: "flags.read"}, {RoleName: "user", PermissionName: "config.read"},
	}

	if err := s.SeedRoles(ctx, roles); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to seed roles: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Seeded roles.")

	if err := s.SeedPermissions(ctx, permissions); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to seed permissions: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Seeded permissions.")

	if err := s.SeedRolePermissions(ctx, mappings); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to seed role-permission mappings: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Seeded role-permission mappings.")

	fmt.Println("Seeding complete.")
}
