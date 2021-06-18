package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"gihtub.com/jimdhughes/stock-watcher/models"
	"gihtub.com/jimdhughes/stock-watcher/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var infosToCheck []models.Watcher
var client http.Client = http.Client{}

var runCommand = &cobra.Command{
	Use:   "run",
	Short: "run monitor",
	Long:  "run monitor for stock as maintained in the application's config.json",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running")
		var err error
		infosToCheck, err = initializeChecks()
		if err != nil {
			log.Fatal(err)
		}
		task := util.Task{
			Closed: make(chan struct{}),
			Ticker: time.NewTicker(time.Second * time.Duration(viper.GetInt32("tick"))),
		}
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		task.Wg.Add(1)
		go func() { defer task.Wg.Done(); task.Run(handleChecks) }()

		sig := <-c
		log.Printf("Got %s signal. Aborting...\n", sig)
	},
}

func init() {
	rootCmd.AddCommand(runCommand)
	runCommand.PersistentFlags().String("config", "../config.json", "config file (default is ../config.json)")
	runCommand.PersistentFlags().Int32("tick", 30, "seconds between ticks for lookup")
	viper.BindPFlag("config", runCommand.PersistentFlags().Lookup("config"))
	viper.BindPFlag("tick", runCommand.PersistentFlags().Lookup("tick"))
}

func initializeChecks() ([]models.Watcher, error) {
	file, err := ioutil.ReadFile(viper.GetString("config"))
	if err != nil {
		return nil, err
	}
	data := []models.Watcher{}
	err = json.Unmarshal([]byte(file), &data)
	return data, err
}

func handleChecks() {
	if len(infosToCheck) < 1 {
		log.Fatal("Nothing to run")
	}
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
			c.HandleSuccess()
		}
		if !strings.Contains(sel, c.LookFor) {
			c.HandleFailure()
		}
	default:
		log.Fatalf("Invalid checktype declared: %s\n", c.CheckType)
	}
}
