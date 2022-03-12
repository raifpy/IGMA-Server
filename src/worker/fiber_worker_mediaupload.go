package worker

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"soccerapi/src/worker/types"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var errJobIdMissing = errors.New("job_id missing")
var errJobIdNotValid = errors.New("job_id not valid")

func (w *Worker) fiberMediaUpload(c *fiber.Ctx) error {

	t := c.Locals("worker").(types.DatabaseAuth)
	jobidstring := c.Query("id", "")
	if jobidstring == "" {
		return errJobIdMissing
	}

	jobid, err := strconv.Atoi(jobidstring)
	if err != nil {
		return errJobIdMissing
	}

	res, err := w.GetJobDb(c.Context(), int64(jobid))
	if err != nil {
		log.Println("w.GetJobDb:fiber_worker_mediaupload.go", err)
		return err //TODO
	}
	if res.WorkerID != t.Id {
		return errJobIdNotValid
	}

	file, err := c.FormFile("media")
	if err != nil {
		log.Println("c.FormFile media error: ", err)
		return err
	}

	fileopen, err := file.Open()
	if err != nil {
		log.Println("file.Open media error: ", err)
		return err
	}
	defer fileopen.Close()

	/*go w.MongoDB.Database.Database("worker").Collection("jobs").FindOneAndReplace(c.Context(), structres.PrivateJobStruct{
		WorkerID: t.Id,
	}, pjobs)*/

	os.Mkdir(path.Join("rendered_videos"), os.ModePerm)
	fullpath := path.Join("rendered_videos", fmt.Sprintf("%s_%d_from%d", strings.Replace(file.Filename, "/", "!converted!", -1), jobid, t.Id))
	writer, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	wint, err := io.Copy(writer, fileopen)
	if err != nil {
		//Todo:
		return err
	}
	fmt.Println("Copied: ", wint)
	fileopen.Close()
	writer.Close()

	var stdout []byte
	if std, err := c.FormFile("stdout"); err == nil {
		if f, err := std.Open(); err == nil {
			defer f.Close()
			stdout, _ = io.ReadAll(f)
		}
	}

	/*go w.MongoDB.Database.Database("worker").Collection("jobs").FindOneAndReplace(context.Background(), structres.PrivateJobStruct{
		WorkerID: t.Id,
	}, pjobs)*/

	res.JobResponseStore = &types.JobResponseStore{
		Size:      wint,
		Local:     true,
		LocalPath: fullpath,

		StdoutSize: int64(len(stdout)),
		Stdout:     stdout,
	} // Match eklemesi yapılacak. Buraya mı?

	if err := w.UpdateJobDb(c.Context(), res); err != nil {
		return err
	}

	return c.SendStatus(200)
}
