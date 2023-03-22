package logs

import (
	"io"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type FormatterHook struct {
	Writer		io.Writer
	LogLevels	[]log.Level
	Formatter	log.Formatter
}

func (hook *FormatterHook) Fire(entry *log.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

func (hook *FormatterHook) Levels() []log.Level {
	return hook.LogLevels
}

func Init() *os.File {
	file, _ := os.OpenFile("./logs/logs.log", os.O_APPEND|os.O_WRONLY, 0600)
	log.SetOutput(ioutil.Discard)
	log.AddHook(&FormatterHook{
		Writer: os.Stderr,
		LogLevels: []log.Level{
			log.DebugLevel,
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
			log.InfoLevel,
		},
		Formatter: &log.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceColors:     true,
		},
	})
	log.AddHook(&FormatterHook{
		Writer: file,
		LogLevels: []log.Level{
			log.DebugLevel,
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
			log.InfoLevel,
		},
		Formatter: &log.JSONFormatter{},
	})
	return file
}

func Elapsed(name string) func() {
	start := time.Now()
	return func() {log.Infof("%s took %v", name, time.Since(start))}
}
