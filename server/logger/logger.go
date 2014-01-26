package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type level int

const (
	TRACE level = 10
	DEBUG level = 20
	INFO  level = 30
	WARN  level = 40
	ERROR level = 50
	FATAL level = 60
)

var (
	pid       int
	name      string        = "-"
	hostname  string        = "-"
	version   int           = 0
	stdout    io.Writer     = os.Stdout
	encoder   *json.Encoder = json.NewEncoder(stdout)
	min_level level         = INFO
)

func init() {
	pid = os.Getpid()
	if v, err := os.Hostname(); err == nil {
		hostname = v
	}
}

type record struct {
	Version  int                    `json:"v"`
	Name     string                 `json:"name,omitempty"`
	Hostname string                 `json:"hostname,omitempty"`
	Pid      int                    `json:"pid,omitempty"`
	Level    level                  `json:"level,omitempty"`
	Created  time.Time              `json:"time"`
	Message  string                 `json:"msg"`
	Fields   map[string]interface{} `json:"fields,omitempty"`
}

func SetName(v string) {
	name = v
}

func SetLevel(v level) {
	min_level = v
}

func Trace(msg string) {
	write(TRACE, msg)
}

func Tracef(msg string, args ...interface{}) {
	write(TRACE, fmt.Sprintf(msg, args...))
}

func Debug(msg string) {
	write(DEBUG, msg)
}

func Debugf(msg string, args ...interface{}) {
	write(DEBUG, fmt.Sprintf(msg, args...))
}

func Info(msg string) {
	write(INFO, msg)
}

func Infof(msg string, args ...interface{}) {
	write(INFO, fmt.Sprintf(msg, args...))
}

func Warn(msg string) {
	write(WARN, msg)
}

func Warnf(msg string, args ...interface{}) {
	write(WARN, fmt.Sprintf(msg, args...))
}

func Error(msg string) {
	write(ERROR, msg)
}

func Errorf(msg string, args ...interface{}) {
	write(ERROR, fmt.Sprintf(msg, args...))
}

func Fatal(msg string) {
	write(FATAL, msg)
}

func Fatalf(msg string, args ...interface{}) {
	write(FATAL, fmt.Sprintf(msg, args...))
}

func write(lvl level, msg string) {
	if lvl < min_level {
		return
	}

	r := record{
		Name:     name,
		Hostname: hostname,
		Pid:      pid,
		Level:    lvl,
		Created:  time.Now(),
		Message:  msg,
	}

	writeRecord(&r)
}

func writeRecord(req *record) {
	req.Version = version

	if req.Name == "" {
		req.Name = name
	}

	if req.Hostname == "" {
		req.Hostname = hostname
	}

	encoder.Encode(req)
}
