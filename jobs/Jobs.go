package jobs

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const JobTypeScraper = "Scraper"
const JobTypeParser = "Parser"

const JobStatusStarted = "Started"
const JobStatusDone = "Done"

type Job struct {
	StartTime, EndTime        int64
	JobId, JobType, JobStatus string
	Params                    map[string]string
}

// Keeps track of running jobs in memory. No need to persist since jobs will only run in
// the single go process for now.
type JobManager struct {
	jobs     map[string]*Job
	jobsLock *sync.RWMutex
}

var jobManagerInstance *JobManager

func GetJobManager() *JobManager {
	if jobManagerInstance == nil {
		jobManagerInstance = &JobManager{}
		jobManagerInstance.jobs = make(map[string]*Job)
		jobManagerInstance.jobsLock = new(sync.RWMutex)
	}

	return jobManagerInstance
}

func (jm *JobManager) GetJobs() map[string]*Job {
	// need a copy to protect against races
	jm.jobsLock.Lock()
	mapCopy := make(map[string]*Job)
	for k, v := range jm.jobs {
		mapCopy[k] = v
	}
	jm.jobsLock.Unlock()
	return mapCopy
}

func (jm *JobManager) AddJob(jobType string, params map[string]string) string {
	newJob := &Job{JobType: jobType, Params: params}

	newJob.StartTime = time.Now().Unix()
	newJob.JobStatus = JobStatusStarted
	newJob.EndTime = 0

	jobId := strconv.Itoa(rand.Int())
	newJob.JobId = jobId

	jm.jobsLock.Lock()
	jm.jobs[jobId] = newJob
	jm.jobsLock.Unlock()

	return jobId
}

func (jm *JobManager) JobComplete(jobId string) {
	jm.jobsLock.Lock()
	job, ok := jm.jobs[jobId]
	if ok {
		job.JobStatus = JobStatusDone
		job.EndTime = time.Now().Unix()
	}
	jm.jobsLock.Unlock()
}

func (jm *JobManager) StartJob(jobFunc func(params map[string]string), jobType string, params map[string]string) string {
	jobId := jm.AddJob(jobType, params)
	jm.jobsLock.Lock()
	go jm.ExecuteJob(jobFunc, jm.jobs[jobId])
	jm.jobsLock.Unlock()

	return jobId
}

func (jm *JobManager) ExecuteJob(jobFunc func(params map[string]string), job *Job) {
	jobFunc(job.Params)
	jm.JobComplete(job.JobId)
}
