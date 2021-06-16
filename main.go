package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

var infosToCheck []CheckInfo
var client http.Client = http.Client{}

const (
	SMTP_HOST     = "SMTP_HOST"
	SMPT_PORT     = "SMTP_PORT"
	SMTP_EMAIL    = "SMTP_EMAIL"
	SMTP_PASSWORD = "SMTP_PASSWORD"
)

func main() {
	initializeEnv()
	InitializeRuntime()
	infos, err := initializeChecks()
	if err != nil {
		log.Fatal("Unable to parse configuration file")
	}
	infosToCheck = infos
	task := &Task{
		closed: make(chan struct{}),
		ticker: time.NewTicker(time.Second * time.Duration(runtimeConfig.tickerDuration)),
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	task.wg.Add(1)
	go func() { defer task.wg.Done(); task.Run() }()

	sig := <-c
	log.Printf("Got %s signal. Aborting...\n", sig)
}

func initializeEnv() {
	godotenv.Load()
	AppMailer = &Mailer{
		SmtpHost: os.Getenv(SMTP_HOST),
		SmtpPort: os.Getenv(SMPT_PORT),
		Email:    os.Getenv(SMTP_EMAIL),
		Password: os.Getenv(SMTP_PASSWORD),
	}
}

func initializeChecks() ([]CheckInfo, error) {
	file, err := ioutil.ReadFile(runtimeConfig.configFileLocation)
	if err != nil {
		return nil, err
	}
	data := []CheckInfo{}
	err = json.Unmarshal([]byte(file), &data)
	return data, err
}

func handleChecks() {
	for _, c := range infosToCheck {
		go handleCheck(c)
	}
}

func handleCheck(c CheckInfo) {
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

func handlePageCheck(doc *goquery.Document, c CheckInfo) {

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
