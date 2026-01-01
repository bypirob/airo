package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"bypirob/airo/src/internal/config"
)

func sshCommand(cfg config.Config, remoteArgs ...string) *exec.Cmd {
	args := []string{}
	if cfg.Deploy.SSH.Port != 0 {
		args = append(args, "-p", fmt.Sprintf("%d", cfg.Deploy.SSH.Port))
	}
	if cfg.Deploy.SSH.IdentityFile != "" {
		args = append(args, "-i", expandUserPath(cfg.Deploy.SSH.IdentityFile))
	}

	sshTarget := cfg.Deploy.SSH.Host
	if cfg.Deploy.SSH.User != "" {
		sshTarget = fmt.Sprintf("%s@%s", cfg.Deploy.SSH.User, cfg.Deploy.SSH.Host)
	}
	args = append(args, sshTarget)
	args = append(args, remoteArgs...)

	return exec.Command("ssh", args...)
}

func expandUserPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, strings.TrimPrefix(path, "~/"))
		}
	}
	return path
}

func shellQuote(value string) string {
	if value == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(value, "'", `'"'"'`) + "'"
}

func shellJoin(args []string) string {
	quoted := make([]string, 0, len(args))
	for _, arg := range args {
		quoted = append(quoted, shellQuote(arg))
	}
	return strings.Join(quoted, " ")
}
