package util

import (
	"fmt"
	"os"
	"time"

	nestedlogrus "github.com/antonfisher/nested-logrus-formatter"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-extras/elogrus.v7"
)

func timestampNow() int64 {
	return time.Now().UnixNano()
}

type LoggerWrapper struct {
	InternalLogger *logrus.Logger
}

func (lw *LoggerWrapper) Initialize() {
	lw.InternalLogger = logrus.New()
	lw.InternalLogger.SetOutput(os.Stdout) // elasticsearch receives from stdout by default
	lw.InternalLogger.SetLevel(logrus.DebugLevel)

	// Formatter
	//log.SetFormatter(&ecslogrus.Formatter{})
	//MyLogger.SetFormatter(&logrus.TextFormatter{
	//	TimestampFormat: "2006-01-02 15:04:05",
	//})

	lw.InternalLogger.SetFormatter(&nestedlogrus.Formatter{
		HideKeys: true,
		//FieldsOrder: []string{"component", "category"},
		NoColors:        true,
		TrimMessages:    true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	lw.InternalLogger.SetReportCaller(true) // log method name by default

	fmt.Println("InternalLogger created")
}

func (lw *LoggerWrapper) Error(args ...interface{}) {
	lw.InternalLogger.WithField("ts", timestampNow()).Error(args...)
}

func (lw *LoggerWrapper) Errorf(format string, args ...interface{}) {
	lw.InternalLogger.WithField("ts", timestampNow()).Errorf(format, args...)
}

func (lw *LoggerWrapper) Warn(args ...interface{}) {
	lw.InternalLogger.WithField("ts", timestampNow()).Warn(args...)
}

func (lw *LoggerWrapper) Warnf(format string, args ...interface{}) {
	lw.InternalLogger.WithField("ts", timestampNow()).Warnf(format, args...)
}

func (lw *LoggerWrapper) Info(args ...interface{}) {
	lw.InternalLogger.WithField("ts", timestampNow()).Info(args...)
}

func (lw *LoggerWrapper) Infof(format string, args ...interface{}) {
	lw.InternalLogger.WithField("ts", timestampNow()).Infof(format, args...)
}

func (lw *LoggerWrapper) Debug(args ...interface{}) {
	lw.InternalLogger.WithField("ts", timestampNow()).Debug(args...)
}

func (lw *LoggerWrapper) Debugf(format string, args ...interface{}) {
	lw.InternalLogger.WithField("ts", timestampNow()).Debugf(format, args...)
}

func (lw *LoggerWrapper) AddElasticHook(indexName string, asyncHook bool) (logrus.Hook, error) {
	cert, err := os.ReadFile("/home/ubuntu/http_ca.crt")
	if err != nil {
		return nil, fmt.Errorf("error reading CA certificate: %s", err)
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
		return nil, err
	}

	// Async hooks mean un-ordered messages dispatch
	// To view logs in kibana Observability Logs, the undex name follows logs-* pattern
	var hook *elogrus.ElasticHook
	if asyncHook {
		hook, err = elogrus.NewAsyncElasticHook(client, "localhost", logrus.DebugLevel, indexName)
		fmt.Printf("ElasticHook async %v\n", asyncHook)
	} else {
		hook, err = elogrus.NewElasticHook(client, "localhost", logrus.DebugLevel, indexName)
		fmt.Printf("ElasticHook async %v\n", asyncHook)
	}

	if err != nil {
		return nil, err
	}
	lw.InternalLogger.Hooks.Add(hook)

	return hook, nil
}
