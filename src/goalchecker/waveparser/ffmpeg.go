package waveparser

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func VideoToWavFFMPEG(path string, outpath string, start, stop string) (string, error) {
	/*format, err := VideoToAudioFormat(path)
	if err != nil {
		return "", err
	}*/
	//outpath = fmt.Sprintf("%s.%s", outpath, strings.ToLower(strings.Split(format, " ")[0]))
	outpath = fmt.Sprintf("%s.%s", outpath, "mp3")
	//cmd := exec.Command("ffmpeg", "-i", path, "-vn", "-acodec", "copy", outpath) //-vn -acodec copy
	cmd := exec.Command("ffmpeg", "-y", "-i", path, "-vn", "-ss", start, "-to", stop, "-v", "quiet", outpath)
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr

	return outpath, cmd.Run()

}

func VideoToAudioFormat(path string) (string, error) {
	out, err := exec.Command("mediainfo", "--inform=Audio;%Format%", path).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), "\n"), nil
}
