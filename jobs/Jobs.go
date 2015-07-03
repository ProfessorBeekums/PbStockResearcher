package jobs

import (
	"math/rand"
	"strconv"
	"time"
)

const JobTypeScreener = "Screener"
const JobTypeParser = "Parser"

const JobStatusStarted = "Started"
const JobStatusDone = "Done"

type Job struct {
	StartTime, EndTime int64
	JobType, JobStatus string
	Params map[string]string
}

// Keeps track of running jobs in memory. No need to persist since jobs will only run in 
// the single go process for now.
type JobManager struct {
	jobs map[string]*Job
}

var jobManagerInstance *JobManager

func GetJobManager() *JobManager {
	if jobManagerInstance == nil {
		jobManagerInstance = &JobManager{}
		jobManagerInstance.jobs = make(map[string]*Job)
	}

	return jobManagerInstance
}

func (jm *JobManager) GetJobs() map[string]*Job {
	return jm.jobs
}

func (jm *JobManager) AddJob(jobType string, params map[string]string) string {
	newJob := &Job{JobType: jobType, Params: params}

	newJob.StartTime = time.Now().Unix()
	newJob.JobStatus = JobStatusStarted
	newJob.EndTime = 0

	jobId := strconv.Itoa(rand.Int())

	jm.jobs[jobId] = newJob

	return jobId
}

func (jm *JobManager) JobComplete(jobId string) {
	job, ok := jm.jobs[jobId]
	if ok {
		job.JobStatus = JobStatusDone
		job.EndTime = time.Now().Unix()
	}
}

/**
TODO notes
web request calls function to execute screener job
screener job creates a channel and runs the screener function in a go routine
screener job creates a job in the job manager
screener function sends data to channel when done.
screener job removes job from the job manager
*/