package dksync

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type CronSync struct {
	// Specifies Concurrency with
	// Available options
	// 	"allow"
	// 	"forbid"
	Concurrency 	string

	// On create allows to enable/disable Job synchronization
	DisableJob  	bool

	// Specifies job Retries on fail
	Retries 		int

	// Specifies Tags
	Tags        	string

	// Specifies Executor type
	// Available options
	//	"shell"
	//	"http"
	//	"rabbitmq"
	Executor 		string

	// Specifies Processor type
	// "mail": {...}
	// "log": {...}
	// "redis": {...}
	Processor		map[string]map[string]string

	// Optional field
	// Specifies Job DisplayName
	// if not empty will be used instead of OS Hostname
	DisplayName    	string

	// test before sync
	DryRun 			bool

	// Unexported fields
	// will be assigned only when fetch crontab file
	cmd 			string
	internalName 	string
	env 			string
}

var (
	hostBuilder strings.Builder
	hostname string
	err error
)

// Processor retrieves jobs from a crontab file
// then sync to API
func (s *CronSync) JobsProcessor(path string) {
	checkHost()
	source, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer source.Close()

	content, err := ioutil.ReadAll(source)
	if err != nil {
		log.Error(err)
	}

	if s.DryRun {
		log.Infoln("Run simulation...")
	}

	// parse crontab file
	crontab := strings.Split(parser(content), "\n")
	for rawCmd := range crontab {
		if crontab[rawCmd] != "" {
			cmd := strings.Fields(crontab[rawCmd])
			if len(cmd) > 4 {
				if s.env = parseEnv(cmd[5]); s.env != "" {
					s.cmd = strings.Join(cmd[6:], " ")
				} else {
					s.cmd = strings.Join(cmd[5:], " ")
				}
				s.scheduler(cmd[:5], rawCmd)
			}
		}
	}
	log.Infoln("Crontab file synchronized")
}

// scheduler build payload
// check if newly defined job already exists
// if not create job
func (s *CronSync) scheduler(cmd []string, id int) {
	s.getHostname(id)

	payload := &job {
		Name:        s.internalName,
		DisplayName: fmt.Sprintf("%s-%s", s.internalName, parseJobName(s.cmd)),
		Timezone:    "Europe/Lisbon",
		Schedule:    fmt.Sprintf("0 %s", strings.Join(cmd, " ")),
		Owner:       "KK",
		OwnerEmail:  "tecnica@kk.pt",
		Disabled:    s.DisableJob,
		Metadata:    nil,
		Processors:  s.Processor,
		Retries:     s.Retries,
		Concurrency: s.Concurrency,
		Executor:    s.Executor,
	}
	payload.ExecutorConfig.Command = s.cmd
	payload.ExecutorConfig.Env = s.env
	if err := json.Unmarshal([]byte(s.Tags), &payload.Tags); err != nil {
		log.Fatal(err)
	}
	if s.Executor == "shell" {
		payload.ExecutorConfig.Shell = "true"
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Errorln(err)
	}
	if status := payload.checkJob(); status == 200 {
		if s.DryRun {
			log.Infoln(string(jsonPayload))
		}
		return
	}

	// simulate sync
	if s.DryRun {
		log.Infoln(string(jsonPayload))
		return
	}
	payload.create()
}

// getHostname get current node hostname
// then replace hostname letter "." with "_"
func (s *CronSync) getHostname(id int) {
	hostname, err = os.Hostname()
	if err != nil {
		if s.DisplayName != "" {
			hostBuilder.WriteString(s.DisplayName)
		} else {
			hostBuilder.WriteString("KK-Node-")
		}
		hostBuilder.WriteString(strconv.Itoa(id))
	} else {
		if s.DisplayName != "" {
			hostBuilder.WriteString(s.DisplayName)
		} else {
			hostBuilder.WriteString(hostname)
		}
		hostBuilder.WriteString("-")
		hostBuilder.WriteString(strconv.Itoa(id))
	}
	s.internalName = strings.ToLower(
		strings.ReplaceAll(hostBuilder.String(), ".", "_"))
}

// parseJobName parses job name from executable file
func parseJobName(job string) string {
	split := regexp.MustCompile("([.][a-z]* )").Split(job, -1)
	rawSplit := strings.Split(split[0], "/")
	return rawSplit[len(rawSplit) - 1]
}

// parseEnv parses environment variables
func parseEnv(command string) string {
	split := strings.Split(command, "=")
	if len(split) > 1 {
		return command
	}
	return ""
}