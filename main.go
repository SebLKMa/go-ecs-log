package main

import (
	"fmt"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
	"gopkg.in/go-extras/elogrus.v7"
)

func console() {
	// Create a new logrus instance
	log := logrus.New()

	// Set the output format to the ECS Formatter
	log.SetFormatter(&ecslogrus.Formatter{})

	// Example log message
	log.Info("Hello, ECS logging in Go!")
}

func hooklog1() error {
	cert, err := os.ReadFile("/home/ubuntu/http_ca.crt")
	if err != nil {
		return fmt.Errorf("error reading CA certificate: %s", err)
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200", // Use https for secure connections
		},
		Username: "elastic",              // Your Elasticsearch username
		Password: "_YHicxB7pLvI-xjMWVVf", // Your Elasticsearch password
		CACert:   cert,                   // Provide the CA certificate bytes here
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return err
	}
	/*
			client, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{"https://localhost:9200"},
			Username:  "elastic",
			Password:  "xSdhglJ_4ohN1IpByxtv",
			CACert:    cert,
		})
	*/

	log := logrus.New()
	log.SetOutput(os.Stdout) // elasticsearch receives from stdout by default
	log.SetFormatter(&ecslogrus.Formatter{})
	log.SetLevel(logrus.DebugLevel)
	hook, err := elogrus.NewAsyncElasticHook(client, "localhost", logrus.DebugLevel, "mylog")
	if err != nil {
		return err
	}
	log.Hooks.Add(hook)
	defer hook.Cancel()

	// Check if you can create views for fields in kibana
	log.WithFields(logrus.Fields{
		"myfield1": "joe",
		"myfield2": 42,
	}).Error("Hello from log hook")

	// anonymous struct
	//msg := struct {
	//	Message   string `json:"message"`
	//	Timestamp int64  `json:"timestamp"`
	//}{}
	msg := struct {
		Message   string
		Timestamp int64
	}{}
	log.Info("Logging started")
	msg.Message = "hello!"
	msg.Timestamp = time.Now().UnixNano()
	log.Debugf("elastic hook log: %#v", msg)
	msg.Message = "bonjour!"
	msg.Timestamp = time.Now().UnixNano()
	log.Infof("elastic hook log: %#v", msg)
	msg.Message = "hola!"
	msg.Timestamp = time.Now().UnixNano()
	log.Warnf("elastic hook log: %#v", msg)
	msg.Message = "oops!"
	msg.Timestamp = time.Now().UnixNano()
	log.Errorf("elastic hook log: %#v", msg)
	log.Info("Logging ended")

	// Give some time for asynchronous logs to be sent
	time.Sleep(2 * time.Second)

	return nil
}

func main() {
	//console()
	err := hooklog1()
	if err != nil {
		panic(err)
	}
}
