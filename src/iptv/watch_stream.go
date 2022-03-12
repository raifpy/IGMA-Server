package iptv

import (
	"context"
	"fmt"
	"io"
	"time"
)

func (ip *IpTv) watchstream(writer io.Writer, channelpath string, onerror func(error)) (int64, error) {
	_url := ip.cloneURL(channelpath)

	var ctx context.Context
	if ip.DeadlineTimeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(ip.DeadlineTimeout))
		defer cancel()
	} else {
		ctx = context.Background()
	}
	request, err := ip.getRequestContext("GET", _url.String(), ctx)
	if err != nil {
		return 0, nil
	}
	var counter = 0
request:
	response, err := ip.Client.Do(request)
	if err != nil {
		if counter < ip.ReRequestOnErrorCount {
			counter++
			time.Sleep(time.Second)
			goto request
		}
		return 0, err
	}
	defer response.Body.Close()
	fmt.Printf("response.Header: %v\n", response.Header)
	if err := ip.ControlResponse(response); err != nil {
		return 0, err
	}

	size, err := io.Copy(writer, response.Body)
	return size, ip.ControlStreamResponseError(err)

}
