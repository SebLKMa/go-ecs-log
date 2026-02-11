package util

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	nestedlogrus "github.com/antonfisher/nested-logrus-formatter"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-extras/elogrus.v7"
)

func timestampNow() int64 {
	return time.Now().UnixNano()
}

// ContextHook adds file and function information to log entries.
type ContextHook struct{}

func (hook ContextHook) Fire(entry *logrus.Entry) error {
	// Skip the logrus call stack frames to find the actual caller
	pc := make([]uintptr, 3, 3)
	cnt := runtime.Callers(6, pc)
	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		// Only use info if it's outside of logrus's own functions
		if !strings.Contains(name, "github.com/sirupsen/logrus") {
			file, line := fu.FileLine(pc[i] - 1)
			entry.Data["file"] = path.Base(file) // Just the filename
			entry.Data["func"] = path.Base(name) // Just the function name
			entry.Data["line"] = line
			break
		}
	}
	return nil
}

func (hook ContextHook) Levels() []logrus.Level {
	// Fire for all levels
	return logrus.AllLevels
}

type LoggerWrapper struct {
	InternalLogger *logrus.Logger
}

func (lw *LoggerWrapper) Initialize() {
	lw.InternalLogger = logrus.New()
	lw.InternalLogger.SetOutput(os.Stdout) // elasticsearch receives from stdout by default
	lw.InternalLogger.SetLevel(logrus.DebugLevel)

	// Formatters
	/*
		logger.SetFormatter(&ecslogrus.Formatter{})
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	*/

	lw.InternalLogger.SetFormatter(&logrus.JSONFormatter{})

	lw.InternalLogger.SetFormatter(&nestedlogrus.Formatter{
		HideKeys: true,
		//FieldsOrder: []string{"component", "category"},
		NoColors:        true,
		TrimMessages:    true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	lw.InternalLogger.SetReportCaller(true) // log method name by default

	//lw.InternalLogger.AddHook(&ContextHook{})

	fmt.Println("InternalLogger created")
}

func (lw *LoggerWrapper) AddElasticHook(indexName string, asyncHook bool) (logrus.Hook, error) {
	cert, err := os.ReadFile("./http_ca.crt")
	if err != nil {
		return nil, fmt.Errorf("error reading CA certificate: %s", err)
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200", // Use https for secure connections
		},
		Username: "elastic",      // Your Elasticsearch username
		Password: "toerrishuman", // Your Elasticsearch password
		CACert:   cert,           // Provide the CA certificate bytes here
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

// Fix to show the actual line calling log
// https://stackoverflow.com/questions/63658002/is-it-possible-to-wrap-logrus-logger-functions-without-losing-the-line-number-pr
/*
func Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Info(args...)
	}
}
*/

func printStack() {
	// Ask runtime.Callers for up to 10 PCs
	pc := make([]uintptr, 10)
	n := runtime.Callers(1, pc) // Skip 1 to ignore printStack itself
	if n == 0 {
		return
	}

	pc = pc[:n] // Pass only valid PCs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	fmt.Println("Call Stack:")
	for {
		frame, more := frames.Next()
		// Customize function name presentation for brevity
		functionName := strings.TrimPrefix(frame.Function, "main.")
		fmt.Printf("* %s: %s#%d\n", functionName, frame.File, frame.Line)

		if !more {
			break
		}
	}
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} /*else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}*/
	return fmt.Sprintf("%s:%d", file, line)
}

func (lw *LoggerWrapper) Error(args ...interface{}) {
	if lw.InternalLogger.Level >= logrus.ErrorLevel {
		entry := lw.InternalLogger.WithFields(logrus.Fields{})
		entry.Data["srcfile"] = fileInfo(2) // file will be updated by logrus, changed file to custom srcfile
		entry.Data["ts"] = timestampNow()

		fmt.Printf("srcfile: %v\n", entry.Data["srcfile"])

		entry.Error(args...)
	}
	//lw.InternalLogger.WithField("ts", timestampNow()).Error(args...)
}

func (lw *LoggerWrapper) Errorf(format string, args ...interface{}) {
	if lw.InternalLogger.Level >= logrus.ErrorLevel {
		entry := lw.InternalLogger.WithFields(logrus.Fields{})
		entry.Data["srcfile"] = fileInfo(2)
		entry.Data["ts"] = timestampNow()

		fmt.Printf("srcfile: %v\n", entry.Data["srcfile"])

		entry.Errorf(format, args...)
	}
	//lw.InternalLogger.WithField("ts", timestampNow()).Errorf(format, args...)
}

func (lw *LoggerWrapper) Warn(args ...interface{}) {
	if lw.InternalLogger.Level >= logrus.WarnLevel {
		entry := lw.InternalLogger.WithFields(logrus.Fields{})
		entry.Data["srcfile"] = fileInfo(2)
		entry.Data["ts"] = timestampNow()
		entry.Warn(args...)
	}
	//lw.InternalLogger.WithField("ts", timestampNow()).Warn(args...)
}

func (lw *LoggerWrapper) Warnf(format string, args ...interface{}) {
	if lw.InternalLogger.Level >= logrus.WarnLevel {
		entry := lw.InternalLogger.WithFields(logrus.Fields{})
		entry.Data["srcfile"] = fileInfo(2)
		entry.Data["ts"] = timestampNow()
		entry.Warnf(format, args...)
	}
	//lw.InternalLogger.WithField("ts", timestampNow()).Warnf(format, args...)
}

func (lw *LoggerWrapper) Info(args ...interface{}) {
	if lw.InternalLogger.Level >= logrus.InfoLevel {
		entry := lw.InternalLogger.WithFields(logrus.Fields{})
		entry.Data["srcfile"] = fileInfo(2)
		entry.Data["ts"] = timestampNow()
		entry.Info(args...)
	}
	//lw.InternalLogger.WithField("ts", timestampNow()).Info(args...)
}

func (lw *LoggerWrapper) Infof(format string, args ...interface{}) {
	if lw.InternalLogger.Level >= logrus.InfoLevel {
		entry := lw.InternalLogger.WithFields(logrus.Fields{})
		entry.Data["srcfile"] = fileInfo(2)
		entry.Data["ts"] = timestampNow()
		entry.Infof(format, args...)
	}
	//lw.InternalLogger.WithField("ts", timestampNow()).Infof(format, args...)
}

func (lw *LoggerWrapper) Debug(args ...interface{}) {
	if lw.InternalLogger.Level >= logrus.DebugLevel {
		entry := lw.InternalLogger.WithFields(logrus.Fields{})
		entry.Data["srcfile"] = fileInfo(2)
		entry.Data["ts"] = timestampNow()
		entry.Debug(args...)
	}
	//lw.InternalLogger.WithField("ts", timestampNow()).Debug(args...)
}

func (lw *LoggerWrapper) Debugf(format string, args ...interface{}) {
	if lw.InternalLogger.Level >= logrus.DebugLevel {
		entry := lw.InternalLogger.WithFields(logrus.Fields{})
		entry.Data["srcfile"] = fileInfo(2)
		entry.Data["ts"] = timestampNow()
		entry.Debugf(format, args...)
	}
	//lw.InternalLogger.WithField("ts", timestampNow()).Debugf(format, args...)
}
