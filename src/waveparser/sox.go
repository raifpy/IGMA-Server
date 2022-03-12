package waveparser

import "os/exec"

func CheckSox() bool {
	_, err := exec.LookPath("sox")
	return err == nil
}
