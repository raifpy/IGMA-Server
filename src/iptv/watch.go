package iptv

import (
	"fmt"
	"io"
	"strings"
	"time"
)

func (i *IpTv) Watch(writer io.Writer, channelpath string, onerror func(error)) (int64, error) {
	//i.Client.Get(i.URL + channelpath)

	fmt.Printf("i.Stream: %v\n", i.Stream)

	if i.Stream {
		return i.watchstream(writer, channelpath, onerror)
	}
	return i.watch(writer, channelpath, onerror)

}

func (ip *IpTv) watch(writer io.Writer, channelpath string, onerror func(error)) (int64, error) {
	var path = channelpath
	var counter int
	var mediapaths = stringList{}
	var size int64

	var subpath string
	cpsplit := strings.Split(channelpath, "/")
	cpslen := len(cpsplit)
	if cpslen > 1 {
		subpath = cpsplit[0] + "/" //!! DÜZELTİLMELİ |=|
	}
requestandcontrol:
	cip := ip.cloneURL(path)
	//log.Println("Requesting ", cip.String())
	response, err := ip.RequestAndControl(cip.String())
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	details, err := parse(response.Body)
	if err != nil {
		return 0, err
	}
	if details.Redirect {
		path = details.RedirectPath
		goto requestandcontrol
	}

	if details.Media {
		for _, p := range details.MediaPaths {
			fmt.Println("details.MediaPaths loop")
			if mediapaths.Contains(p) {
				continue
			}

			var requrl string = p

			if !details.UsesOtherURL {
				curl := ip.cloneURL(subpath + p)

				requrl = curl.String()
			}
			fmt.Printf("p octet/stream url: %v\n", requrl)

			response, err = ip.RequestAndControl(requrl)
			if err != nil {
				onerror(err)
				continue
			}
			//fmt.Printf("response.Header: %v\n", response.Header)
			defer response.Body.Close()

			_size, err := io.Copy(writer, response.Body)
			fmt.Println(_size, err)
			if err != nil {
				//ip.OnError(err)
				onerror(err)
				continue
			}
			size += _size
			mediapaths = append(mediapaths, p)
		}
		counter++
		fmt.Println(counter, ip.Loop)
		if counter < ip.Loop {
			time.Sleep(time.Second * time.Duration(ip.LoopSleep))
			goto requestandcontrol
		}

		//log.Println("Tamam galiba")

	}
	return size, nil
}

/*
func (ip *IpTv) watchonotherurl(_url string, writer io.Writer) (int64, error) {
	log.Println("watchonotherurl: ", _url)
	var size int64
	response, err := ip.RequestAndControl(_url)
	if err != nil {
		return size, err
	}
	defer response.Body.Close()
	details, err := parse(response.Body)
	if err != nil {
		return size, err
	}
}
*/
