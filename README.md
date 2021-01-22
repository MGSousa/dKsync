# dKsync

Package tools for import / export cron jobs using [dKron API](https://github.com/distribworks/dkron) 

## Install
```shell script
go get -v github.com/MGSousa/dKsync
```

## Run

### Use Cron Job Importer
```go
package main

import "github.com/MGSousa/dKsync"

func main() {
	sync := dksync.CronSync {
		Concurrency:    "allow",
		DisableJob:     false,
		Tags:           "{\"mytag\":\"true\"}",
		Executor:       "shell",
		Processor: 	    map[string]map[string]string{"plugin": {...}},
		Retries:      	1,
		DryRun: 	    true,
	}
	sync.JobsProcessor("/var/spool/cron/crontabs/root")
}
```