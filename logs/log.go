package logs

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	color_gray = uint8(iota + 90)
	color_red
	color_green
	color_yellow
	color_blue
	color_magenta //洋红
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var defaultLogger = NewHLogger()

type HLogger struct {
	core   *log.Logger
	levelS []string
	//输出的最低日志等级
	targetLevel  LogLevel
	lvlMaxLength int
}

/**
基于标准Log库封装的，支持日志等级，彩色显示的日志库
默认最低日志等级
*/
func NewHLogger() *HLogger {
	l := log.New(os.Stderr, "", log.Ltime|log.Ldate)
	g := &HLogger{
		core: l,
	}
	levelS := make([]string, 5)
	levelS[LevelDebug] = gray("DEBUG")
	levelS[LevelInfo] = blue("INFO")
	levelS[LevelWarn] = yellow("WARN")
	levelS[LevelError] = red("ERROR")
	levelS[LevelFatal] = magenta("FATAL")
	maxLen := len(levelS[0])
	for _, lvl := range levelS {
		if len(lvl) > maxLen {
			maxLen = len(lvl)
		}
	}
	g.lvlMaxLength = maxLen
	g.levelS = levelS
	g.targetLevel = LevelDebug
	return g
}

func (g *HLogger) SetLevel(i LogLevel) {
	g.targetLevel = i
}

func (g *HLogger) Debug(s ...interface{}) {
	g.log(LevelDebug, s...)
}

func (g *HLogger) Info(s ...interface{}) {
	g.log(LevelInfo, s...)
}

func (g *HLogger) Warn(s ...interface{}) {
	g.log(LevelWarn, s...)
}

func (g *HLogger) Error(s ...interface{}) {
	g.log(LevelError, s...)
}

func (g *HLogger) Fatal(s ...interface{}) {
	g.log(LevelFatal, s...)
	msgs := fmt.Sprintln(s...)
	panic(strings.TrimSuffix(msgs, "\n"))
}

/**
封装方法
*/
func SetLevel(i LogLevel) {
	defaultLogger.targetLevel = i
}

func Debug(s ...interface{}) {
	defaultLogger.Debug(s...)
}

func Info(s ...interface{}) {
	defaultLogger.Info(s...)
}

func Warn(s ...interface{}) {
	defaultLogger.Warn(s...)
}

func Error(s ...interface{}) {
	defaultLogger.Error(s...)
}

func Fatal(s ...interface{}) {
	defaultLogger.Fatal(s...)
	msgs := fmt.Sprintln(s...)
	panic(strings.TrimSuffix(msgs, "\n"))
}

/**
core function
*/
func (g *HLogger) log(lvl LogLevel, s ...interface{}) {
	//omit unconcerned logs
	if lvl < g.targetLevel {
		return
	}
	fileName, lineNumber := "", 0
	skip := 2
	if g == defaultLogger {
		skip = 3
	}
	pc, _, _, ok := runtime.Caller(skip)
	if ok {
		fileName, lineNumber = runtime.FuncForPC(pc).FileLine(pc)
		for i := len(fileName) - 1; i > 0; i-- {
			if fileName[i] == '/' {
				fileName = fileName[i+1:]
				break
			}
		}
	}
	if ok {
		apiName := runtime.FuncForPC(pc).Name()
		methods := strings.Split(apiName, ".")
		if len(methods) > 1 {
			apiName = methods[len(methods)-1]
		}
	}
	lvlS := g.getAlignedLevel(lvl)
	//unify all messages
	msgs := fmt.Sprintln(s...)
	msgs = strings.TrimSuffix(msgs, "\n")
	//align levels

	g.core.Printf("%s %s:%d %s", lvlS, fileName, lineNumber, msgs)
}

func (g *HLogger) getAlignedLevel(i LogLevel) string {
	lvl := g.levelS[i]
	if c := g.lvlMaxLength - len(lvl); c > 0 {
		for i := 0; i < c; i++ {
			lvl += " "
		}
	}
	return lvl
}

func gray(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_gray, s)
}

func red(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_red, s)
}

func green(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_green, s)
}

func yellow(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_yellow, s)
}

func blue(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_blue, s)
}

func magenta(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_magenta, s)
}
