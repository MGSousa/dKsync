package dksync

import (
	"fmt"
	"github.com/gemalto/requester"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type (
	job struct {
		Name           string                 		`json:"name"`
		DisplayName    string                 		`json:"displayname"`
		Timezone       string                 		`json:"timezone"`
		Schedule       string                 		`json:"schedule"`
		Owner          string                 		`json:"owner"`
		OwnerEmail     string                 		`json:"owner_email"`
		Disabled       bool                   		`json:"disabled"`
		Tags           map[string]string 			`json:"tags"`
		Metadata       interface{}            		`json:"metadata"`
		Processors     map[string]map[string]string	`json:"processors"`
		Retries        int                    		`json:"retries"`
		Concurrency    string                 		`json:"concurrency"`
		Executor       string                 		`json:"executor"`
		ExecutorConfig struct {
			Command    string 	`json:"command"`
			Cwd 	   string 	`json:"cwd,omitempty"`
			Env 	   string 	`json:"env,omitempty"`
			Shell 	   string 	`json:"shell"`
		} `json:"executor_config"`
	}
)

// create
func (job *job) create() {
	resp, err := requester.Send(
		requester.JSON(false),
		requester.Body(job),
		requester.Post(dKronHost),
		requester.ExpectCode(201),
		requester.AddHeader("Accept", "application/json"),
		requester.AddHeader("Content-Type", "application/json"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Debugf("New Job: %v", resp)

	log.Infoln(resp.Status)
	return
}

// update
func (job *job) update() {
	resp, err := requester.Send(
		requester.Post(dKronHost),
		requester.Body(&job),
		requester.ExpectCode(201),
		requester.AddHeader("Accept", "application/json"),
		requester.AddHeader("Content-Type", "application/json"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Debugf("Updated Job: %s, Name: %s", resp.Status, job.Name)
	return
}

// checkJob check if job already exists
func (job *job) checkJob() int {
	response, err := requester.Send(
		requester.Get(fmt.Sprintf("%s%s", dKronHost, job.Name)))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	log.Debugln(response)
	return response.StatusCode
}

// getJobs
func getJobs() []byte {
	resp, err := requester.Send(
		requester.Get(dKronHost),
		requester.ExpectCode(200),
		requester.AddHeader("Accept", "application/json"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Debugf("Get Jobs: %v", resp)

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

// checkHost
func checkHost() {
	if dKronHost == "" {
		dKronHost = defaultHost
	}
}

// setDebug
func setDebug() {
	log.SetLevel(log.DebugLevel)
}