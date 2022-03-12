package types

type JobDB struct {
	Job              Job               `json:"job,omitempty" mongo:"job,omitempty"`
	JobId            int64             `json:"job_id,omitempty" mongo:"job_id,omitempty"`
	WorkerID         int64             `json:"worker_id,omitempty" mongo:"worker_id,omitempty"`
	LastUpdates      []Job             `json:"last_updates,omitempty" mongo:"last_updates,omitempty"`
	JobResponseStore *JobResponseStore `json:"job_response_store,omitempty" mongo:"job_response_store,omitempty"`
}

type JobResponseStore struct {
	Size int64 `json:"size"`

	Local     bool   `json:"local"`
	LocalPath string `json:"local_path"`

	Gdrive   bool   `json:"gdrive"`
	GdriveId string `json:"gdrive_id"`

	StdoutSize int64  `json:"stdout_size"`
	Stdout     []byte `json:"stdout"`
}
