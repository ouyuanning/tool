package github

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type newHostContent struct {
	content    string
	lastUpdate string
}
type oldHostContent struct {
	content    []string
	hostStart  int
	hostEnd    int
	lastUpdate string
}

func getNewHost(hostUrl string) *newHostContent {
	resp, err := http.Get(hostUrl)
	if err != nil {
		log.Println("println日志222sss")
		log.Printf("%v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("println日志222dddddd")
		log.Printf("%v", err)
		os.Exit(1)
	}
	newContent := string(body)
	var lastUpdate string
	for _, line := range strings.Split(newContent, "\n") {
		if strings.HasPrefix(line, "# Update time:") {
			lastUpdate = strings.SplitN(line, ":", 2)[1]
		}
	}
	return &newHostContent{
		content:    newContent,
		lastUpdate: lastUpdate,
	}
}

func getOldContent(hostFilePath string) *oldHostContent {
	content, err := ioutil.ReadFile(hostFilePath)
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}
	oldContent := strings.Split(string(content), "\n")
	var hostStart int
	var hostFinish int
	var lastUpdate string
	for idx, line := range oldContent {
		if line == "# GitHub520 Host Start" {
			hostStart = idx
		}
		if line == "# GitHub520 Host End" {
			hostFinish = idx
		}
		if strings.HasPrefix(line, "# Update time:") {
			lastUpdate = strings.SplitN(line, ":", 2)[1]
		}
	}
	return &oldHostContent{
		content:    oldContent,
		hostStart:  hostStart,
		hostEnd:    hostFinish,
		lastUpdate: lastUpdate,
	}
}

func ScanHostAndSetToLocal() {
	hostFilePath := "/etc/hosts"
	oldContent := getOldContent(hostFilePath)

	hostUrl := "https://raw.hellogithub.com/hosts"
	newContent := getNewHost(hostUrl)

	if oldContent.lastUpdate != newContent.lastUpdate {
		// log.Printf("%v, %v", oldContent.lastUpdate, newContent.lastUpdate)
		updateContent := ""
		idx := 0
		for idx < oldContent.hostStart {
			updateContent = updateContent + oldContent.content[idx] + "\n"
			idx++
		}
		updateContent = updateContent + newContent.content
		idx = oldContent.hostEnd + 1
		max := len(oldContent.content) - 1
		for idx <= max {
			updateContent = updateContent + oldContent.content[idx]
			if idx < max {
				updateContent = updateContent + "\n"
			}
			idx++
		}

		err := ioutil.WriteFile(hostFilePath, []byte(updateContent), 0644)
		if err != nil {
			log.Printf("%v", err)
			os.Exit(1)
		}
		log.Printf("change hosts ok")
	} else {
		log.Println("no change")
	}
}
