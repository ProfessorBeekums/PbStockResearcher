package jobs

import (
	"sync"
	"testing"
	"time"
)

func TestJobMap(t *testing.T) {
	jm := GetJobManager()

	jobMap := jm.GetJobs()

	// test initial state
	if len(jobMap) > 0 {
		t.Fatal("Job map is not empty")
	}

	params := make(map[string]string)
	params["param1"] = "123"
	params["param2"] = "456"

	jobId1 := jm.AddJob("blar", params)
	
	// verify singleton
	jm = GetJobManager()

	jobMap = jm.GetJobs()

	// verify added job data
	if len(jobMap) != 1 {
		t.Fatal("Job map does not have one entry")
	}

	v := jobMap[jobId1]
	if v.JobType != "blar" {
		t.Fatal("Incorrect job type. Excepted blar, got: ", v.JobType)
	}

	if len(v.Params) != 2 {
		t.Fatal("Incorrect number of params. Expected 2, got: ", len(v.Params))
	}

	jm.AddJob("blar", params)
	
	jobMap = jm.GetJobs()

	// verify multiple jobs
	if len(jobMap) != 2 {
		t.Fatal("Job map does not have expected 2 entries. Received: ", len(jobMap))
	}

	jm.JobComplete(jobId1)

	jobMap = jm.GetJobs()

	completedJob := jobMap[jobId1]

	if completedJob.JobStatus != JobStatusDone {
		t.Fatal("Job <", jobId1, "> was expected to have done status, instead got: ", completedJob.JobStatus)
	}

	if completedJob.EndTime == 0 {
		t.Fatal("Job <", jobId1, "> was expected to have end time of 0")
	}

	if completedJob.JobId != jobId1 {
		t.Fatal("Mismatched job id <", jobId1, "> and : ", completedJob.JobId)
	}
}

func TestStartJob(t *testing.T) {
	jm := GetJobManager()

	jobMap := jm.GetJobs()

	params := make(map[string]string)
	params["param1"] = "123"

	jobId := jm.StartJob(fakeJobFunc, "blar", params)

	// not the best, but the ops are guaranteed fast since it's just setting a variable
	time.Sleep(5 * time.Millisecond)

	fakeParamLock.Lock()
	if fakeParam != "123" {
		t.Fatal("Job function was no executed by StartJob")
	}
	fakeParamLock.Unlock()

	jobMap = jm.GetJobs()
	completedJob := jobMap[jobId]

	if completedJob.JobStatus != JobStatusDone {
		t.Fatal("JobComplete was not executed by StartJob")
	}
}

var fakeParamLock *sync.RWMutex = new(sync.RWMutex)
var fakeParam string

func fakeJobFunc(params map[string]string) {
	fakeParamLock.Lock()
	fakeParam = params["param1"]
	fakeParamLock.Unlock()
}