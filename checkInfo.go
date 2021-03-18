package main

import (
	"log"

	"github.com/gookit/color"
)

type CheckInfo struct {
	URL              string `json:"url"`
	Key              string `json:"key"`
	Description      string `json:"description"`
	LookFor          string `json:"lookFor"`
	OnSuccessMessage string `json:"onSuccessMessage"`
	OnFailureMessage string `json:"onFailureMessage"`
	MailTo           string `json:"mailTo"`
	SmsTo            string `json:"smsTo"`
	CheckType        string `json:"checkType"`
	IsNegativeCheck  bool   `json:"isNegativeCheck"`
	Vendor           string `json:"vendor"`
}

func (c *CheckInfo) HandleLogEvent(success bool) {
	if c.IsNegativeCheck == true {
		success = !success
	}
	if success {
		log.Printf("[%s] %s : %s @ %s\n", c.Vendor, c.Key, color.FgGreen.Render(c.OnSuccessMessage), c.URL)
	}
	if !success {
		log.Printf("[%s] %s : %s @ %s\n", c.Vendor, c.Key, color.FgRed.Render(c.OnFailureMessage), c.URL)
	}
}
