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

	//log.WithFields(logrus.Fields{
	//	"name": "joe",
	//	"age":  42,
	//}).Error("Hello from log hook")

	log.Debug("Debug ECS logging from elastic hook")
	log.Info("Info ECS logging from elastic hook")
	log.Warn("Warning ECS logging from elastic hook")
	log.Error("Error ECS logging from elastic hook")

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
