package waveparser

type FrameInfos []FrameInfo

func (fs FrameInfos) GetBiggestUpFrequency() (vü FrameInfo) {
	for i, abcd := range fs {
		if abcd.UpWave > vü.UpWave {
			vü = abcd
			vü.Index = i
		}
	}
	return
}

func (fs FrameInfos) GetBiggestDownFrequency() (vü FrameInfo) {
	for i, abcd := range fs {
		if abcd.DownWave > vü.DownWave {
			vü = abcd
			vü.Index = i
		}
	}
	return
}
