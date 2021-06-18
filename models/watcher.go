package models

import (
	"fmt"
	"log"
	"time"

	"gihtub.com/jimdhughes/stock-watcher/util"
	"github.com/gookit/color"
)

type CustomHeaders struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Watcher struct {
	URL              string          `json:"url"`
	Key              string          `json:"key"`
	Description      string          `json:"description"`
	LookFor          string          `json:"lookFor"`
	OnSuccessMessage string          `json:"onSuccessMessage"`
	OnFailureMessage string          `json:"onFailureMessage"`
	MailTo           []string        `json:"mailTo"`
	SmsTo            string          `json:"smsTo"`
	CheckType        string          `json:"checkType"`
	IsNegativeCheck  bool            `json:"isNegativeCheck"`
	Vendor           string          `json:"vendor"`
	CustomHeaders    []CustomHeaders `json:"headers"`
}

func (c *Watcher) HandleLogEvent(success bool) {
	if success {
		log.Printf("[%s] %s : %s @ %s\n", c.Vendor, c.Key, color.FgGreen.Render(c.OnSuccessMessage), c.URL)
	}
	if !success {
		log.Printf("[%s] %s : %s @ %s\n", c.Vendor, c.Key, color.FgRed.Render(c.OnFailureMessage), c.URL)
	}
}

func (c *Watcher) GetMailMessage() string {
	datetime := time.Now()
	msg := fmt.Sprintf("[%s] %s : %s @ %s\n\n Detected At: %s\n", c.Vendor, c.Key, c.OnSuccessMessage, c.URL, datetime.Format("2006-01-02 15:04:05"))
	return msg
}

func (c *Watcher) GetMailSubject() string {
	return fmt.Sprintf("[%s] %s is %s", c.Vendor, c.Key, c.OnSuccessMessage)
}

func (c *Watcher) HandleMail(success bool) {
	if !success {
		return
	}
	err := util.AppMailer.SendMail(c.MailTo, c.GetMailMessage(), c.GetMailSubject())
	if err != nil {
		log.Printf("ERROR trying to send mail: %s\n", err.Error())
	}
}

func (c *Watcher) HandleFailure() {
	success := false
	if c.IsNegativeCheck {
		success = true
	}
	c.HandleLogEvent(success)
	c.HandleMail(success)
}

func (c *Watcher) HandleSuccess() {
	success := true
	if c.IsNegativeCheck {
		success = false
	}
	c.HandleLogEvent(success)
	c.HandleMail(success)
}
