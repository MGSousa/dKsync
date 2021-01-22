package dksync

const defaultHost = "http://127.0.0.1:8080/v1/jobs/"

type Opts struct {
	// allows crontab file synchronization
	CronSync

	// allows API synchronization updates
	ApiSync
}

// SetHost for dKron API
func SetHost(url string) {
	dKronHost = url
}
