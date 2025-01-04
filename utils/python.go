package utils

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	matchesLen = 4
)

type PythonVersion struct {
	Major int
	Minor int
	Patch int
}

func GetPythonVersion() (*PythonVersion, error) {
	cmd := exec.Command("python3", "-V")

	output, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute command\n")
	}

	re := regexp.MustCompile(`Python (\d+)\.(\d+)\.(\d+)`)

	matches := re.FindStringSubmatch(strings.TrimSpace(string(output)))
	if len(matches) != matchesLen {
		return nil, errors.New("failed to parse version\n")
	}

	version := &PythonVersion{}

	if _, err := fmt.Sscanf(matches[1], "%d", &version.Major); err != nil {
		return nil, errors.Wrap(err, "failed to parse major version\n")
	}

	if _, err := fmt.Sscanf(matches[2], "%d", &version.Minor); err != nil {
		return nil, errors.Wrap(err, "failed to parse minor version\n")
	}

	if _, err := fmt.Sscanf(matches[3], "%d", &version.Patch); err != nil {
		return nil, errors.Wrap(err, "failed to parse patch version\n")
	}

	return version, nil
}
