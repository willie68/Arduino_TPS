package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aphistic/golf"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	LvDebug string = "DEBUG"
	LvInfo  string = "INFO"
	LvAlert string = "ALERT"
	LvError string = "ERROR"
	LvFatal string = "FATAL"
)

var Levels = []string{LvDebug, LvInfo, LvAlert, LvError, LvFatal}

/*
ServiceLogger main type for logging
*/
type serviceLogger struct {
	Level      string
	LevelInt   int
	GelfURL    string
	GelfPort   int
	SystemID   string
	Attrs      map[string]interface{}
	gelfActive bool
	c          *golf.Client
	Filename   string
}

// Logger to use for all logging
var Logger serviceLogger

/*
Debug log this message at debug level
*/
func Debug(msg string) {
	Logger.Debug(msg)
}

/*
Debugf log this message at debug level with formatting
*/
func Debugf(format string, va ...interface{}) {
	Logger.Debugf(format, va...)
}

/*
Info log this message at info level
*/
func Info(msg string) {
	Logger.Info(msg)
}

/*
Infof log this message at info level with formatting
*/
func Infof(format string, va ...interface{}) {
	Logger.Infof(format, va...)
}

/*
Alert log this message at alert level
*/
func Alert(msg string) {
	Logger.Alert(msg)
}

/*
Alertf log this message at alert level with formatting.
*/
func Alertf(format string, va ...interface{}) {
	Logger.Alertf(format, va...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(msg string) {
	Logger.Fatal(msg)
}

// Fatalf logs a message at level Fatal on the standard logger with formatting.
func Fatalf(format string, va ...interface{}) {
	Logger.Fatalf(format, va...)
}

// Error logs a message at level Error on the standard logger.
func Error(msg string) {
	Logger.Error(msg)
}

// Errorf logs a message at level Error on the standard logger with formatting.
func Errorf(format string, va ...interface{}) {
	Logger.Errorf(format, va...)
}

/*
Init initialise logging
*/
func (s *serviceLogger) Init() {
	s.gelfActive = false
	if s.GelfURL != "" {
		s.c, _ = golf.NewClient()
		s.c.Dial(fmt.Sprintf("udp://%s:%d", s.GelfURL, s.GelfPort))

		l, _ := s.c.NewLogger()

		golf.DefaultLogger(l)
		for key, value := range s.Attrs {
			l.SetAttr(key, value)
		}
		l.SetAttr("system_id", s.SystemID)
		s.gelfActive = true
	}
	var w io.Writer
	if s.Filename == "" {
		w = os.Stdout
	} else {
		w = io.MultiWriter(&lumberjack.Logger{
			Filename:   s.Filename,
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28,    //days
			Compress:   false, // disabled by default
		}, os.Stdout)
	}
	log.SetOutput(w)
}

func (s *serviceLogger) SetLevel(level string) {
	switch strings.ToUpper(level) {
	case LvDebug:
		s.LevelInt = 0
	case LvInfo:
		s.LevelInt = 1
	case LvAlert:
		s.LevelInt = 2
	case LvError:
		s.LevelInt = 3
	case LvFatal:
		s.LevelInt = 4
	}
}

/*
Debug log this message at debug level
*/
func (s *serviceLogger) Debug(msg string) {
	if s.LevelInt <= 0 {
		if s.gelfActive {
			golf.Dbg(msg)
		}
		log.Printf("Debug: %s\n", msg)
	}
}

/*
Debugf log this message at debug level with formatting
*/
func (s *serviceLogger) Debugf(format string, va ...interface{}) {
	if s.LevelInt <= 0 {
		if s.gelfActive {
			golf.Dbgf(format, va...)
		}
		log.Printf("Debug: %s\n", fmt.Sprintf(format, va...))
	}
}

/*
Info log this message at info level
*/
func (s *serviceLogger) Info(msg string) {
	if s.LevelInt <= 1 {
		if s.gelfActive {
			golf.Info(msg)
		}
		log.Printf("Info: %s\n", msg)
	}
}

/*
Infof log this message at info level with formatting
*/
func (s *serviceLogger) Infof(format string, va ...interface{}) {
	if s.LevelInt <= 1 {
		if s.gelfActive {
			golf.Infof(format, va...)
		}
		log.Printf("Info: %s\n", fmt.Sprintf(format, va...))
	}
}

/*
Alert log this message at alert level
*/
func (s *serviceLogger) Alert(msg string) {
	if s.LevelInt <= 2 {
		if s.gelfActive {
			golf.Alert(msg)
		}
		log.Printf("Alert: %s\n", msg)
	}
}

/*
Alertf log this message at alert level with formatting.
*/
func (s *serviceLogger) Alertf(format string, va ...interface{}) {
	if s.LevelInt <= 2 {
		if s.gelfActive {
			golf.Alertf(format, va...)
		}
		log.Printf("Alert: %s\n", fmt.Sprintf(format, va...))
	}
}

// Fatal logs a message at level Fatal on the standard logger.
func (s *serviceLogger) Fatal(msg string) {
	if s.LevelInt <= 4 {
		if s.gelfActive {
			golf.Crit(msg)
		}
		log.Fatalf("Fatal: %s\n", msg)
	}
}

// Fatalf logs a message at level Fatal on the standard logger with formatting.
func (s *serviceLogger) Fatalf(format string, va ...interface{}) {
	if s.LevelInt <= 4 {
		if s.gelfActive {
			golf.Critf(format, va...)
		}
		log.Fatalf("Fatal: %s\n", fmt.Sprintf(format, va...))
	}
}

// Error logs a message at level Error on the standard logger.
func (s *serviceLogger) Error(msg string) {
	if s.LevelInt <= 3 {
		if s.gelfActive {
			golf.Err(msg)
		}
		log.Printf("Error: %s\n", msg)
	}
}

// Errorf logs a message at level Error on the standard logger with formatting.
func (s *serviceLogger) Errorf(format string, va ...interface{}) {
	if s.LevelInt <= 3 {
		if s.gelfActive {
			golf.Errf(format, va...)
		}
		log.Printf("Error: %s\n", fmt.Sprintf(format, va...))
	}
}

/*
Close this logging client
*/
func (s *serviceLogger) Close() {
	if s.gelfActive {
		s.c.Close()
	}
}
