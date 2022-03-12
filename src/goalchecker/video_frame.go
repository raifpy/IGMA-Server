package goalchecker

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
)

func GetFrameFromVideo(sec string, videopath string, output string) error {

	os.Remove(output)

	if err := exec.Command("ffmpeg", "-y", "-v", "quiet", "-ss", sec, "-i", videopath, "-vframes", "1", "-vf", "crop=450:300:0:0", output).Run(); err != nil {
		return err
	}
	_, err := os.Stat(output)
	return err
}

func TryGetLastFrameFromVideo(lastsec int, videopath string, output string) (int, error) {

	for {
		if lastsec <= 0 {
			break
		}
		if err := GetFrameFromVideo(strconv.Itoa(lastsec), videopath, output); err == nil {
			return lastsec, nil
		}
		lastsec--

	}

	return 0, errors.New("unexcepted error")

}
