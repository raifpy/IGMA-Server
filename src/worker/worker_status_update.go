package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"soccerapi/src/worker/types"
)

//!!! Buraya update şart.

func (w *Worker) WorkerStatusUpdate(wt types.WebsocketContact, conn *WorkerConn) {
	defer func() {
		rec := recover() // :')
		if rec != nil {
			if res, _ := http.Get("https://api.telegram.org/bot<TOKEN>/sendMessage?chat_id=<TELEGRAM_ID>&text=" + fmt.Sprint(rec)); res != nil {
				res.Body.Close()
			}
		}
	}()

	r, err := json.MarshalIndent(wt, " ", " ")
	if err != nil {
		log.Println("WorkerStatusUpdate jsonMarshalIndent: ", err)
		return

	}

	fmt.Println(string(r))

	// Belki doğrudan buraya chan trigger

	switch wt.Type {
	case "close":
		log.Printf("%d close gönderdi", conn.Id)
		w.WorkerMap.Del(conn)

		conn.CloseConnection()
		conn.Conn.Close()

		if wt.Error != nil {
			log.Println(*wt.Error)
			return
		}

	case "done", "uploading", "downloading", "rendering":
		if wt.Update != nil {
			go func() {
				res, err := w.GetJobDb(context.Background(), wt.Update.Job.JobID)
				if err != nil {
					//TODO Trigger chan??
					log.Printf("%d's Update request mongo GetJobDb \033[31merror\033[0m: %v", conn.Id, err)
					return
				}
				res.Job.Status = wt.Type

				if len(res.LastUpdates) == 0 {
					res.LastUpdates = []types.Job{res.Job}
				}

				if err := w.UpdateJobDb(context.Background(), res); err != nil {
					log.Printf("w.UpdateJobDb error: %v", err)
				}

			}()

		}

	case "error":
		if wt.Error != nil {
			if wt.Error.Job == nil {
				log.Printf("Worker %d error: \033[31m %s \033[0m", conn.Id, wt.Error.Error)
				//wt.Error.Job.JobID
				return
			}
			log.Printf("Worker %d Job: %d error: \033[31m %s \033[0m", conn.Id, wt.Error.Job.JobID, wt.Error.Error)

			go func() {
				res, err := w.GetJobDb(context.Background(), wt.Error.Job.JobID)
				if err != nil {
					//TODO Trigger chan??
					log.Printf("%d's Update request (error request) mongo GetJobDb \033[31merror\033[0m: %v", conn.Id, err)
					return
				}
				res.Job.Status = wt.Type
				res.Job.Error = wt.Error.Error

				if len(res.LastUpdates) == 0 {
					res.LastUpdates = []types.Job{res.Job}
				}

				if err := w.UpdateJobDb(context.Background(), res); err != nil {
					log.Printf("w.UpdateJobDb error: %v", err)
				}
				log.Println("Error update edildi")

			}()

		}
	}
	var jobid int64

	if wt.Update != nil {
		jobid = wt.Update.Job.JobID
	}
	if wt.Error != nil && wt.Error.Job != nil {
		jobid = wt.Error.Job.JobID
	}

	fmt.Printf("jobid: %v\n", jobid)

	if jobid != 0 {
		if c := w.WorkerMap.GetChan(jobid); c != nil {
			c <- wt
		} else {
			log.Printf("GetChan (jobid %d) empty!", jobid)
		}
	}

	/*
		switch {
		case wt.Error != nil:
			{
				log.Printf("\033[31mERROR\033[0m: %v Client: %d\n", wt.Error.Error, conn.Id)
				if wt.Error.Job != nil {
					go w.UpdateJobDb(context.Background(), types.JobDB{
						Job:      *wt.Error.Job,
						JobId:    wt.Error.Job.JobID,
						WorkerID: conn.Id,
					})
					if c := w.WorkerMap.GetChan(wt.Error.Job.JobID); c != nil {
						c <- wt
						return
					}
				}
			}
		case wt.Update != nil:
			{

				switch wt.Update.Job.Status {

				case "error":
					{
						log.Printf("%d's job update request error: %v", conn.Id, err)
						go w.DelJobDb(context.Background(), wt.Update.Job.ToJobDb())
						return
					}
				case "uploading":
					{
						log.Printf("%d's job update media uploading", conn.Id)
						return
					}
				case "done":
					log.Printf("%d's job update request done!", conn.Id)
					//?? onDone chan trigger? //şart

					return
				}

			}
		}*/

}
