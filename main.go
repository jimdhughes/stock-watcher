package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gihtub.com/jimdhughes/stock-watcher/cmd"
	"gihtub.com/jimdhughes/stock-watcher/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

var infosToCheck []models.Watcher
var client http.Client = http.Client{}

func main() {
	initializeEnv()
	InitializeRuntime()
	cmd := cmd.GetRootCommand()
	cmd.Execute()
}

func initializeEnv() {
	godotenv.Load()
}

func initializeChecks() ([]models.Watcher, error) {
	file, err := ioutil.ReadFile(runtimeConfig.configFileLocation)
	if err != nil {
		return nil, err
	}
	data := []models.Watcher{}
	err = json.Unmarshal([]byte(file), &data)
	return data, err
}

func handleChecks() {
	for _, c := range infosToCheck {
		go handleCheck(c)
	}
}

func handleCheck(c models.Watcher) {
	req, err := http.NewRequest(http.MethodGet, c.URL, nil)
	if len(c.CustomHeaders) > 0 {
		for _, h := range c.CustomHeaders {
			req.Header.Add(h.Key, h.Value)
		}
	}
	if err != nil {
		log.Printf("ERROR forming GET request: %s\n", err.Error())
		return
	}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR Getting URL : %s. ERROR: %s", c.URL, err.Error())
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		log.Printf("HTTP Error code received: %d. Will Retry on next run\n", res.StatusCode)
		return
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	handlePageCheck(doc, c)
}

func handlePageCheck(doc *goquery.Document, c models.Watcher) {

	switch c.CheckType {
	case "className":
		sel := doc.Find(c.LookFor)
		if len(sel.Nodes) > 0 {
			c.HandleSuccess()
		}
		if len(sel.Nodes) == 0 {
			c.HandleFailure()
		}

	case "text":
		sel := doc.Text()
		if strings.Contains(sel, c.LookFor) {
			log.Printf("Found %s on %s\n", c.LookFor, c.Key)
			c.HandleSuccess()
		}
		if !strings.Contains(sel, c.LookFor) {
			c.HandleFailure()
		}
	default:
		log.Fatalf("Invalid checktype declared: %s\n", c.CheckType)
	}
}
