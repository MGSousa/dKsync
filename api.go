package dksync

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type ApiSync struct {
	// Tag name
	Tag  	string

	// Enables Debug
	Debug 	bool

	// Args to update
	Args 	Args

	// Unexported fields
	// jobs
	jobs 	[]job
}

// Processor retrieves jobs from API considering specific query filter
// then update each job in API
func (sync *ApiSync) Processor() {
	checkHost()
	if sync.Debug {
		setDebug()
	}
	if jobs := getJobs(); jobs != nil {
		if err := json.Unmarshal(jobs, &sync.jobs); err != nil {
			log.Fatal(err)
		}
		for _, job := range sync.jobs {
			// TODO - allow to filter for other options
			if job.Tags[sync.Tag] != "" {
				sync.Args.setFieldValue(&job)
				job.update()
			}
		}
	}
}
