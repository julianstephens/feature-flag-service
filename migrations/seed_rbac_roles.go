package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var defaultRoles = []struct {
	Name         string
	Description  string
	PublicRoleID string
}{
	{"admin", "Full access to all resources and administrative functions.", uuid.New().String()},
	{"editor", "Can modify feature flags and configurations.", uuid.New().String()},
	{"auditor", "Read-only access to audit logs and system status.", uuid.New().String()},
	{"user", "Can read feature flags and config, but not modify.", uuid.New().String()},
}

func main() {
	dbURL := os.Getenv("POSTGRES_URL")
	if dbURL == "" {
		fmt.Fprintln(os.Stderr, "POSTGRES_URL env var is required")
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	for _, role := range defaultRoles {
		var exists bool
		err := conn.QueryRow(
			context.Background(),
			`SELECT EXISTS(SELECT 1 FROM rbac_roles WHERE name=$1)`, role.Name,
		).Scan(&exists)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to check role existence: %v\n", err)
			os.Exit(1)
		}
		if !exists {
			_, err = conn.Exec(
				context.Background(),
				`INSERT INTO rbac_roles (name, description, public_role_id)
				VALUES ($1, $2, $3)`,
				role.Name, role.Description, role.PublicRoleID,
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to insert role %q: %v\n", role.Name, err)
				os.Exit(1)
			}
			fmt.Printf("Inserted role: %s (public_role_id: %s)\n", role.Name, role.PublicRoleID)
		} else {
			fmt.Printf("Role already exists: %s\n", role.Name)
		}
	}
	fmt.Println("Seeding complete.")
}