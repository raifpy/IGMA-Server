package soccer

import (
	"log"
	"sort"
)

func (s *Soccer) FilterIpTvResponseList(list []IpTvResponse) (i IpTvResponse, err error) {

	defer func() {
		//!! del others
	}()
	if len(list) == 1 {
		i = list[0]
	}
	sort.Slice(list, func(i, j int) bool {
		return len(list[i].ErrorList) <= len(list[j].ErrorList)
	}) // küçükten büyüğe?

	for _, eleme := range list {
		if eleme.FatalError != nil {
			log.Printf("%s fatal: \033[31m%v\033[0m", eleme.IpTvName, eleme.FatalError)
			continue
		}
		i = eleme // nasıl bir filtre yapacağım bilmiyorum :D
		break
	}
	if i.FatalError != nil {
		return i, nil
	}

	// Buradan sonrasında hata var!
	if len(list) == 1 {
		return i, i.FatalError
	}

	return list[0], list[0].FatalError

	//defer remove(diğerleri?)

	//TODO!!!!!!!!
}

/*
func GetVideoFrequency(videopath string) (waveparser.FrameInfos, error) {
	out := "/tmp/video2_" + RandStringRunes(10) + "_c"
	path, err := waveparser.VideoToWavFFMPEG(videopath, out)
	defer os.Remove(path)
	log.Println("\033[31m" + out + "\033[0m")
	if err != nil {

		return nil, err
	}
	out2 := out + ".dat"
	w := waveparser.WaveParser{
		Input:     path,
		Frequency: "1",
		Output:    out2,
		Stdout:    nil,
	}
	defer os.Remove(out2)
	return w.Parse()
}
*/
