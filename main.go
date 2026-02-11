package main

import (
	"fmt"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/seblkma/go-ecs-log/util"
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
	"gopkg.in/go-extras/elogrus.v7"
)

var MyLogger util.LoggerWrapper

func init() {
	// Moved to loggerwrapper
	/*
		MyLogger = logrus.New()
		MyLogger.SetOutput(os.Stdout)
		MyLogger.SetLevel(logrus.DebugLevel)

		MyLogger.SetFormatter(&nestedlogrus.Formatter{
			HideKeys: true,
			NoColors:        true,
			TrimMessages:    true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
		MyLogger.SetReportCaller(true)
	*/

	MyLogger.Initialize()

	fmt.Println("MyLogger created")
}

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
		Password: "uCl8kHO51qymS79WPzNK", // Your Elasticsearch password
		CACert:   cert,                   // Provide the CA certificate bytes here
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return err
	}

	log := logrus.New()
	log.SetOutput(os.Stdout) // elasticsearch receives from stdout by default
	log.SetLevel(logrus.DebugLevel)

	// formatter
	//log.SetFormatter(&ecslogrus.Formatter{})
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// To view logs in kibana Observability Logs, the undex name follows logs-* pattern
	hook, err := elogrus.NewAsyncElasticHook(client, "localhost", logrus.DebugLevel, "logs-mylog")
	if err != nil {
		return err
	}
	log.Hooks.Add(hook)
	defer hook.Cancel()

	// Check if you can create views for fields in kibana
	/*
		log.WithFields(logrus.Fields{
			"myfield1": "joe",
			"myfield2": 42,
		}).Error("Hello from log hook")
	*/

	// anonymous struct to test logging an object
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

func testFunc3() {
	MyLogger.Info("Entered")
	defer MyLogger.Info("Exit")
}

func testFunc2() {
	MyLogger.Info("Entered")
	defer MyLogger.Info("Exit")
	testFunc3()
}

func testFunc1() {
	MyLogger.Info("Entered")
	defer MyLogger.Info("Exit")
	// anonymous struct to test logging an object
	msg := struct {
		Message   string
		Timestamp int64
	}{}

	msg.Message = "hello!"
	msg.Timestamp = time.Now().UnixNano()
	MyLogger.Debugf("MyLogger: %#v", msg)
	msg.Message = "bonjour!"
	msg.Timestamp = time.Now().UnixNano()
	MyLogger.Infof("MyLogger: %#v", msg)
	msg.Message = "hola!"
	msg.Timestamp = time.Now().UnixNano()
	MyLogger.Warnf("MyLogger: %#v", msg)
	msg.Message = "oops!"
	msg.Timestamp = time.Now().UnixNano()
	MyLogger.Errorf("MyLogger: %#v", msg)

	testFunc2()
}

func main() {
	//console()

	// Simple logrus hook test
	/*
		err := hooklog1()
		if err != nil {
			panic(err)
		}
	*/

	// Add hook to logger
	_, err := MyLogger.AddElasticHook("logs-mylog", true)
	if err != nil {
		panic(err)
	}

	testFunc1()

	// Without logger wrapper
	/*
		MyLogger.WithField("ts", timestampNow()).Info("Logging started")
		msg.Message = "hello!"
		msg.Timestamp = time.Now().UnixNano()
		MyLogger.WithField("ts", timestampNow()).Debugf("MyLogger: %#v", msg)
		msg.Message = "bonjour!"
		msg.Timestamp = time.Now().UnixNano()
		MyLogger.WithField("ts", timestampNow()).Infof("MyLogger: %#v", msg)
		msg.Message = "hola!"
		msg.Timestamp = time.Now().UnixNano()
		MyLogger.WithField("ts", timestampNow()).Warnf("MyLogger: %#v", msg)
		msg.Message = "oops!"
		msg.Timestamp = time.Now().UnixNano()
		MyLogger.WithField("ts", timestampNow()).Errorf("MyLogger: %#v", msg)
		MyLogger.WithField("ts", timestampNow()).Info("Logging ended")
	*/

	// Common logged static fields
	//myLogEntry := MyLogger.WithFields(logrus.Fields{"fieldname": "fieldvalue"})
	/*
		myLogEntry.Info("Logging started")
		msg.Message = "hello!"
		msg.Timestamp = time.Now().UnixNano()
		myLogEntry.Debugf("MyLogger: %#v", msg)
		msg.Message = "bonjour!"
		msg.Timestamp = time.Now().UnixNano()
		myLogEntry.Infof("MyLogger: %#v", msg)
		msg.Message = "hola!"
		msg.Timestamp = time.Now().UnixNano()
		myLogEntry.Warnf("MyLogger: %#v", msg)
		msg.Message = "oops!"
		msg.Timestamp = time.Now().UnixNano()
		myLogEntry.Errorf("MyLogger: %#v", msg)
		myLogEntry.Info("Logging ended")
	*/

	// Give some time for asynchronous logs to be sent
	time.Sleep(2 * time.Second)
}
