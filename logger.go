package logger

import (
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"regexp"
	"runtime"
	"time"
)

const (
	Debug = level(iota)
	Info  = level(iota)
	Warn  = level(iota)
	Error = level(iota)
)

type level int

const (
	NoColor = color("\x1b[0m")
	Red     = color("\x1b[91m")
	Green   = color("\x1b[32m")
	Yellow  = color("\x1b[33m")
	Blue    = color("\x1b[34m")
	Magenta = color("\x1b[35m")
	Cyan    = color("\x1b[36m")
	White   = color("\x1b[37m")
)

type color string

var colors = []color{
	Red,
	Green,
	Yellow,
	Blue,
	Magenta,
	Cyan,
	White,
}

var pattern *regexp.Regexp

func init() {
	pattern = regexp.MustCompile(`\(.*\)\.\w*.*?`)
}

type option func(*logger)

func New(tag string, opts ...option) Logger {
	color := colors[int(hash(tag))%len(colors)]

	logger := &logger{
		tag,
		color,
		Info,
		os.Stdout,
	}

	for _, opt := range opts {
		opt(logger)
	}

	return logger
}

func Writer(w io.Writer) option {
	return func(log *logger) {
		log.writer = w
	}
}

func Color(c color) option {
	return func(log *logger) {
		log.color = c
	}
}

func Level(l level) option {
	return func(log *logger) {
		log.level = l
	}
}

type Logger interface {
	Fatal(a ...interface{})
	Fatalf(format string, a ...interface{})

	Error(a ...interface{})
	Errorf(format string, a ...interface{})

	Warn(a ...interface{})
	Warnf(format string, a ...interface{})

	Info(a ...interface{})
	Infof(format string, a ...interface{})

	Debug(a ...interface{})
	Debugf(format string, a ...interface{})
}

type logger struct {
	tag    string
	color  color
	level  level
	writer io.Writer
}

func (self *logger) Fatal(a ...interface{}) {
	self.log("FATAL", self.format(a))
	os.Exit(1)
}

func (self *logger) Fatalf(f string, a ...interface{}) {
	self.log("FATAL", self.formatf(f, a))
	os.Exit(1)
}

func (self *logger) Error(a ...interface{}) {
	if self.level <= Error {
		self.log("ERROR", self.format(a))
	}
}

func (self *logger) Errorf(f string, a ...interface{}) {
	if self.level <= Error {
		self.log("ERROR", self.formatf(f, a))
	}
}

func (self *logger) Warn(a ...interface{}) {
	if self.level <= Warn {
		self.log("WARN", self.format(a))
	}
}

func (self *logger) Warnf(f string, a ...interface{}) {
	if self.level <= Warn {
		self.log("WARN", self.formatf(f, a))
	}
}

func (self *logger) Info(a ...interface{}) {
	if self.level <= Info {
		self.log("INFO", self.format(a))
	}
}

func (self *logger) Infof(f string, a ...interface{}) {
	if self.level <= Info {
		self.log("INFO", self.formatf(f, a))
	}
}

func (self *logger) Debug(a ...interface{}) {
	if self.level <= Debug {
		self.log("DEBUG", self.format(a))
	}
}

func (self *logger) Debugf(f string, a ...interface{}) {
	if self.level <= Debug {
		self.log("DEBUG", self.formatf(f, a))
	}
}

func (self *logger) format(a []interface{}) string {
	return fmt.Sprint(a...)
}

func (self *logger) formatf(f string, a []interface{}) string {
	return fmt.Sprintf(f, a...)
}

func (self *logger) log(level string, s string) {
	if self.tag != "" {
		tag := fmt.Sprintf("%s%s%s", self.color, self.tag, NoColor)
		fmt.Fprintf(self.writer, "[%s] [%5s] [%s] [%s] %s\n", timestamp(), level, tag, caller(), s)
	} else {
		fmt.Fprintf(self.writer, "[%s] [%5s] [%s] %s\n", timestamp(), level, caller(), s)
	}
}

func timestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func caller() string {
	pc, _, _, ok := runtime.Caller(3)
	if !ok {
		return "unknown caller"
	}

	details := runtime.FuncForPC(pc)
	if details == nil {
		return "unknown caller"
	}

	match := pattern.FindString(details.Name())
	if match != "" {
		return match
	}

	return details.Name()
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
