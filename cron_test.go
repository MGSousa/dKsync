package dksync

import (
	"testing"
)

type args struct {
	path string
}

func TestCronSync_JobsProcessor(t *testing.T) {
	tests := []struct {
		name   string
		args   args
	}{
		{
			name: "Cron1",
			args: args {
				path: "/var/spool/cron/crontabs/root",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CronSync {
				DisableJob:   	true,
				Concurrency:  	"allow",
				Tags:         	"{\"mytag\":\"true\"}",
				Executor:     	"shell",
				Processor: 	  	map[string]map[string]string{
					"mail": {"subject": "Test Mail"},
				},
				Retries:      	1,
				DryRun: 		true,
			}
			s.JobsProcessor(tt.args.path)
		})
	}
}
