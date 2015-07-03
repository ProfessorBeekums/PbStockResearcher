package jobs

import (
	"testing"
)

func TestJobManager(t *testing.T) {
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
}