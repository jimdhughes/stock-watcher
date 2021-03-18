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
)

var infosToCheck []CheckInfo

func main() {
	err, infos := initializeChecks()
	if err != nil {
		log.Fatal("Unable to parse input")
	}
	infosToCheck = infos
	task := &Task{
		closed: make(chan struct{}),
		ticker: time.NewTicker(time.Second * 20),
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	task.wg.Add(1)
	go func() { defer task.wg.Done(); task.Run() }()

	select {
	case sig := <-c:
		log.Printf("Got %s signal. Aborting...\n", sig)
		task.Stop()
	}
}

func initializeChecks() (error, []CheckInfo) {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err, nil
	}
	data := []CheckInfo{}
	err = json.Unmarshal([]byte(file), &data)
	return err, data
}

func handleChecks() {
	for _, c := range infosToCheck {
		go handleCheck(c)
	}
}

func handleCheck(c CheckInfo) {
	res, err := http.Get(c.URL)
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
		sel := doc.Find(".product-out-of-stock")
		if len(sel.Nodes) > 0 {
			c.HandleLogEvent(true)
		}
		if len(sel.Nodes) == 0 {
			c.HandleLogEvent(false)
		}
		break

	case "text":
		sel := doc.Text()
		if strings.Contains(sel, c.LookFor) {
			c.HandleLogEvent(true)
		}
		c.HandleLogEvent(false)
		break
	default:
		log.Fatalf("Invalid checktype declared: %s\n", c.CheckType)
	}

}
