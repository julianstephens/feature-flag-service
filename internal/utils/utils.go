package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/julianstephens/go-utils/cliutil"
)

const DEFAULT_TIMEOUT = 30 * time.Second

func GenerateID() string {
	return uuid.New().String()
}

func PrintUser(id, email, name, createdAt, updatedAt string, roles []string) {
	cliutil.PrintInfo("User Details:")
	cliutil.PrintInfo("-------------")
	if id != "" && id != "N/A" {
		cliutil.PrintInfo(fmt.Sprintf("  ID:        %s", id))
	}
	if email != "" && email != "N/A" {
		cliutil.PrintInfo(fmt.Sprintf("  Email:     %s", email))
	}
	if name != "" && name != "N/A" {
		cliutil.PrintInfo(fmt.Sprintf("  Name:      %s", name))
	}
	if createdAt != "" && createdAt != "N/A" {
		cliutil.PrintInfo(fmt.Sprintf("  Created At: %s", createdAt))
	}
	if updatedAt != "" && updatedAt != "N/A" {
		cliutil.PrintInfo(fmt.Sprintf("  Updated At: %s", updatedAt))
	}
	if len(roles) > 0 {
		cliutil.PrintInfo(fmt.Sprintf("  Roles:     %s", strings.Join(roles, ", ")))
	} else {
		cliutil.PrintInfo("  Roles:     None")
	}
}
