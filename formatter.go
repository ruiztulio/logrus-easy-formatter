// Package easy allows to easily format output of Logrus logger
package easy

import (
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// Default log format will output [INFO]: 2006-01-02T15:04:05Z07:00 - Log message
	defaultLogFormat           = "[%lvl%]: %time% - %msg%"
	defaultTimestampFormat     = time.RFC3339
	knownGocoreVFrames     int = 9
)

// Formatter implements logrus.Formatter interface.
type Formatter struct {
	// Timestamp format
	TimestampFormat string
	// Available standard keys: time, msg, lvl
	// Also can include custom fields but limited to strings.
	// All of fields need to be wrapped inside %% i.e %time% %msg%
	LogFormat string
}

func getCaller() (string, string, int) {
	pc, file, line, ok := runtime.Caller(knownGocoreVFrames)
	if !ok {
		panic("Could not get context info for logger!")
	}
	funcname := runtime.FuncForPC(pc).Name()
	// if the func name container the package name, lets get the next frame
	if strings.Contains(funcname, "gocorev/logger") {
		pc, file, line, ok = runtime.Caller(knownGocoreVFrames+1)
		funcname = runtime.FuncForPC(pc).Name()
	}

	filename := file[strings.LastIndex(file, "/")+1:]
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	return filename, fn, line
}

// Format building log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.LogFormat
	if output == "" {
		output = defaultLogFormat
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)

	output = strings.Replace(output, "%msg%", entry.Message, 1)

	level := strings.ToUpper(entry.Level.String())
	output = strings.Replace(output, "%lvl%", level, 1)
	fname, fn, ln := getCaller()
	output = strings.Replace(output, "%caller%", fname, 1)
	output = strings.Replace(output, "%line%", strconv.Itoa(ln), 1)
	output = strings.Replace(output, "%func%", fn, 1)

	for k, val := range entry.Data {
		switch v := val.(type) {
		case string:
			output = strings.Replace(output, "%"+k+"%", v, 1)
		case int:
			s := strconv.Itoa(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		case bool:
			s := strconv.FormatBool(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		}
	}

	return []byte(output), nil
}
