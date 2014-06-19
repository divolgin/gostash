package gostash

import (
	"encoding/json"
	"github.com/howbazaar/loggo"
	"path/filepath"
	"strconv"
	"time"
)

type Formatter struct {
	message  logstashMessage
	metadata map[string]string
}

type logstashMessage struct {
	Timestamp string            `json:"@timestamp"`
	Message   string            `json:"message"`
	Type      string            `json:"type"`
	Fields    map[string]string `json:"fields"`
}

func NewFormatter(config *Config) *Formatter {
	formatter := Formatter{}

	formatter.metadata = make(map[string]string)
	if config.ThisHostName != "" {
		formatter.metadata["host"] = config.ThisHostName
	}
	if config.CodeVersion != "" {
		formatter.metadata["version"] = config.CodeVersion
	}

	if config.InputType == "" {
		formatter.message.Type = "go"
	} else {
		formatter.message.Type = config.InputType
	}

	return &formatter
}

func (formatter *Formatter) Format(level loggo.Level, module, filename string, line int, timestamp time.Time, message string) string {
	formatter.metadata["level"] = level.String()
	formatter.metadata["module"] = module
	formatter.metadata["file"] = filepath.Base(filename)
	formatter.metadata["line"] = strconv.Itoa(line)
	formatter.metadata["log.time"] = timestamp.In(time.UTC).Format("2006-01-02 15:04:05")
	encoded, err := formatter.encode(message, formatter.metadata)
	if err != nil {
		return ""
	}
	return string(encoded[:])
}

func (formatter *Formatter) encode(msg string, metadata map[string]string) ([]byte, error) {
	formatter.message.Timestamp = getTimeStamp()
	formatter.message.Message = msg
	formatter.message.Fields = metadata
	return json.Marshal(formatter.message)
}

func getTimeStamp() string {
	now := time.Now().UTC()
	return now.Format("2006-01-02T15:04:05.000-07:00")
}
