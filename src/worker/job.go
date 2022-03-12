package worker

import (
	"context"
	"fmt"
	"math/rand"
	"soccerapi/src/worker/types"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func (w *Worker) AddJobDb(c context.Context, conn *WorkerConn, job types.Job) (int64, error) {
	job.JobID = rand.Int63n(999999999999998)
	fmt.Printf("job.go job.JobID: %v\n", job.JobID)

	if job.Exec != nil {
		var newexec = []string{}
		for _, v := range job.Exec.Args {
			v = strings.Replace(v, "{jobid}", fmt.Sprint(job.JobID), -1)
			newexec = append(newexec, v)
		}

		job.Exec.Args = newexec
	}

	job.Status = "newjob"
	if err := w.SetJobDb(c, types.JobDB{
		JobId:    job.JobID,
		WorkerID: conn.Id,
		Job:      job,
	}); err != nil {
		return job.JobID, err
	}

	err := conn.Conn.WriteJSON(types.WebsocketContact{
		Type:   "newjob",
		NewJob: &job,
	})
	if err != nil {
		go w.DelJobDb(context.Background(), types.JobDB{JobId: job.JobID, WorkerID: conn.Id})
		return job.JobID, err
	}

	return job.JobID, err
}

func (w *Worker) GetJobDb(c context.Context, id int64) (tjb types.JobDB, err error) {
	err = w.MongoDB.Database.Database("jobs").Collection("job").FindOne(c, bson.M{
		"jobid": id,
	}).Decode(&tjb)

	return
}

func (w *Worker) DelJobDb(c context.Context, job types.JobDB) error {

	_, err := w.MongoDB.Database.Database("jobs").Collection("job").DeleteOne(c, types.JobDB{
		Job: types.Job{JobID: job.Job.JobID},
	})
	return err
}

func (w *Worker) SetJobDb(c context.Context, job types.JobDB) error {
	_, err := w.MongoDB.Database.Database("jobs").Collection("job").InsertOne(c, job)
	return err
}

func (w *Worker) UpdateJobDb(c context.Context, job types.JobDB) error {
	return w.MongoDB.Database.Database("jobs").Collection("job").FindOneAndReplace(c, bson.M{
		"workerid": job.WorkerID,
		"jobid":    job.JobId,
	}, job).Err()
}
