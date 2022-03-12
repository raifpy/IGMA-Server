package goalchecker

import (
	"errors"
	"os/exec"
	"strings"
)

func GetVideoDuration(path string) (string, error) {
	out, err := exec.Command("ffprobe", "-i", path, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0").CombinedOutput()
	if err != nil {
		return "", errors.New(string(out))
	}
	return strings.TrimRight(string(out), "\n"), nil
}
