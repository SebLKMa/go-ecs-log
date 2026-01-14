package main

import (
	"log"
	"os"

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

func hook() error {
	cert, err := os.ReadFile("/home/ubuntu/http_ca.crt")
	if err != nil {
		log.Fatalf("Error reading CA certificate: %s", err)
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200", // Use https for secure connections
		},
		Username: "elastic",              // Your Elasticsearch username
		Password: "xSdhglJ_4ohN1IpByxtv", // Your Elasticsearch password
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
	log.SetFormatter(&ecslogrus.Formatter{})
	hook, err := elogrus.NewAsyncElasticHook(client, "localhost", logrus.DebugLevel, "mylog")
	if err != nil {
		return err
	}
	log.Hooks.Add(hook)
	//log.WithFields(logrus.Fields{
	//	"name": "joe",
	//	"age":  42,
	//}).Error("Hello from log hook")

	log.Info("Hello info, ECS logging from elastic hook")
	log.Error("Hello error, ECS logging from elastic hook")

	return nil
}

func main() {
	console()
	err := hook()
	if err != nil {
		panic(err)
	}
}
