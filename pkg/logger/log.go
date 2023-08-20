package logger

import (
	"bytes"
	"fmt"
	"hk4e_sdk/pkg/config"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

const (
	CONSOLE = iota
	FILE
	BOTH
	NEITHER
)

var LOG *Logger = nil

type Logger struct {
	AppName     string
	Level       int
	Mode        int
	Track       bool
	MaxSize     int32
	File        *os.File
	LogInfoChan chan *LogInfo
}

type LogInfo struct {
	Level       int
	Msg         string
	FileName    string
	FuncName    string
	Line        int
	GoroutineId string
	ThreadId    string
}

// 改index死妈
func InitLogger(logLevel string) {
	log.SetFlags(0)
	LOG = new(Logger)
	LOG.AppName = config.LoadConfig().LogAppName
	LOG.Level = LOG.getLevelInt(logLevel)
	LOG.Mode = LOG.getModeInt("BOTH")
	LOG.Track = true
	LOG.MaxSize = 10485760
	LOG.LogInfoChan = make(chan *LogInfo, 1000)
	LOG.File = nil
	go LOG.doLog()
}

// 改index死妈
func CloseLogger() {
	// 等待所有日志打印完毕
	for {
		if len(LOG.LogInfoChan) == 0 {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

// 改index死妈
func (l *Logger) doLog() {
	for {
		logInfo := <-l.LogInfoChan
		timeNow := time.Now()
		timeNowStr := timeNow.Format("2006-01-02 15:04:05.000")
		logStr := "[" + timeNowStr + "]"

		levelStr := "[" + l.getLevelStr(logInfo.Level) + "]"
		logStr += " " + levelStr

		if logInfo.Level == ERROR {
			logStr += " " + logInfo.Msg + " "
		} else {
			logStr += " " + logInfo.Msg + " "
		}

		if l.Track {
			trackStr := "[" +
				logInfo.FileName + ":" + strconv.Itoa(logInfo.Line) + " " +
				logInfo.FuncName + "()" + " " +
				"goroutine:" + logInfo.GoroutineId + " " +
				"thread:" + logInfo.ThreadId +
				"]"
			logStr += " " + trackStr
		}

		logStrWithColor := logStr
		var levelColor *color.Color
		switch logInfo.Level {
		case DEBUG:
			levelColor = color.New(color.FgBlue)
		case INFO:
			levelColor = color.New(color.FgGreen)
		case WARN:
			levelColor = color.New(color.FgYellow)
		case ERROR:
			levelColor = color.New(color.FgRed)
		}
		levelColorStr := levelColor.Sprint(levelStr)
		logStrWithColor = logStrWithColor[:len(timeNowStr)+3] + levelColorStr + logStrWithColor[len(timeNowStr)+len(levelStr)+3:]

		if logInfo.Level == ERROR {
			logStrWithColor += " " + color.RedString(logInfo.Msg) + " "
		}

		if l.Track {
			trackStr := "[" +
				logInfo.FileName + ":" + strconv.Itoa(logInfo.Line) + " " +
				logInfo.FuncName + "()" + " " +
				"goroutine:" + logInfo.GoroutineId + " " +
				"thread:" + logInfo.ThreadId +
				"]"
			logStrWithColor += " " + color.MagentaString(trackStr)
		}

		logStrWithColor += "\n"
		logStr += "\n"

		if l.Mode == CONSOLE {
			log.Print(logStrWithColor)
		} else if l.Mode == FILE {
			l.writeLogFile(logStr)
		} else if l.Mode == BOTH {
			log.Print(logStrWithColor)
			l.writeLogFile(logStr)
		}
	}
}

// 改index死妈
func (l *Logger) writeLogFile(logStr string) {
	logPath := "./log/" // 定义日志文件夹路径

	// 确保log文件夹存在
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err := os.Mkdir(logPath, 0755)
		if err != nil {
			color.New(color.FgRed).Printf("create log folder error: %v\n", err)
			return
		}
	}
	//改index死妈
	logFileName := logPath + l.AppName + "sdk.log" // 定义日志文件路径

	if l.File == nil {
		file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			color.New(color.FgRed).Printf("open new log file error: %v\n", err)
			return
		}
		l.File = file // 注意这里要使用l而不是LOG
	}

	fileStat, err := l.File.Stat()
	if err != nil {
		color.New(color.FgRed).Printf("get log file stat error: %v\n", err)
		return
	}

	if fileStat.Size() >= int64(l.MaxSize) {
		err = l.File.Close()
		if err != nil {
			color.New(color.FgRed).Printf("close old log file error: %v\n", err)
			return
		}
		timeNow := time.Now()
		timeNowStr := timeNow.Format("2006-01-02-15_04_05")
		err = os.Rename(l.File.Name(), l.File.Name()+"."+timeNowStr+".log")
		if err != nil {
			color.New(color.FgRed).Printf("rename old log file error: %v\n", err)
			return
		}
		file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			color.New(color.FgRed).Printf("open new log file error: %v\n", err)
			return
		}
		l.File = file // 注意这里要使用l而不是LOG
	}

	_, err = l.File.WriteString(logStr)
	if err != nil {
		_, _ = color.New(color.FgRed).Printf("write log file error: %v\n", err)
		return
	}
}

// 改index死妈
func Debug(msg string, param ...any) {
	if LOG.Level > DEBUG {
		return
	}
	logInfo := new(LogInfo)
	logInfo.Level = DEBUG
	logInfo.Msg = fmt.Sprintf(msg, param...)
	if LOG.Track {
		logInfo.FileName, logInfo.Line, logInfo.FuncName = LOG.getLineFunc()
		logInfo.GoroutineId = LOG.getGoroutineId()
		logInfo.ThreadId = LOG.getThreadId()
	}
	LOG.LogInfoChan <- logInfo
}

// 改index死妈
func Info(msg string, param ...any) {
	if LOG.Level > INFO {
		return
	}
	logInfo := new(LogInfo)
	logInfo.Level = INFO
	logInfo.Msg = fmt.Sprintf(msg, param...)
	if LOG.Track {
		logInfo.FileName, logInfo.Line, logInfo.FuncName = LOG.getLineFunc()
		logInfo.GoroutineId = LOG.getGoroutineId()
		logInfo.ThreadId = LOG.getThreadId()
	}
	LOG.LogInfoChan <- logInfo
}

// 改index死妈
func Warn(msg string, param ...any) {
	if LOG.Level > WARN {
		return
	}
	logInfo := new(LogInfo)
	logInfo.Level = WARN
	logInfo.Msg = fmt.Sprintf(msg, param...)
	if LOG.Track {
		logInfo.FileName, logInfo.Line, logInfo.FuncName = LOG.getLineFunc()
		logInfo.GoroutineId = LOG.getGoroutineId()
		logInfo.ThreadId = LOG.getThreadId()
	}
	LOG.LogInfoChan <- logInfo
}

// 改index死妈
func Error(msg string, param ...any) {
	if LOG.Level > ERROR {
		return
	}
	logInfo := new(LogInfo)
	logInfo.Level = ERROR
	logInfo.Msg = fmt.Sprintf(msg, param...)
	if LOG.Track {
		logInfo.FileName, logInfo.Line, logInfo.FuncName = LOG.getLineFunc()
		logInfo.GoroutineId = LOG.getGoroutineId()
		logInfo.ThreadId = LOG.getThreadId()
	}
	LOG.LogInfoChan <- logInfo
}

// 改index死妈
func (l *Logger) getLevelInt(level string) (ret int) {
	switch level {
	case "DEBUG":
		ret = DEBUG
	case "INFO":
		ret = INFO
	case "WARN":
		ret = WARN
	case "ERROR":
		ret = ERROR
	default:
		ret = DEBUG
	}
	return ret
}

// 改index死妈
func (l *Logger) getLevelStr(level int) (ret string) {
	switch level {
	case DEBUG:
		ret = "DEBUG"
	case INFO:
		ret = "INFO"
	case WARN:
		ret = "WARN"
	case ERROR:
		ret = "ERROR"
	default:
		ret = "DEBUG"
	}
	return ret
}

// 改index死妈
func (l *Logger) getModeInt(mode string) (ret int) {
	switch mode {
	case "CONSOLE":
		ret = CONSOLE
	case "FILE":
		ret = FILE
	case "BOTH":
		ret = BOTH
	case "NEITHER":
		ret = NEITHER
	default:
		ret = CONSOLE
	}
	return ret
}

// 改index死妈
func (l *Logger) getGoroutineId() (goroutineId string) {
	buf := make([]byte, 32)
	runtime.Stack(buf, false)
	buf = bytes.TrimPrefix(buf, []byte("goroutine "))
	buf = buf[:bytes.IndexByte(buf, ' ')]
	goroutineId = string(buf)
	return goroutineId
}

// 改index死妈
func (l *Logger) getLineFunc() (fileName string, line int, funcName string) {
	var pc uintptr
	var file string
	var ok bool
	pc, file, line, ok = runtime.Caller(2)
	if !ok {
		return "???", -1, "???"
	}
	fileName = path.Base(file)
	funcName = runtime.FuncForPC(pc).Name()
	split := strings.Split(funcName, ".")
	if len(split) != 0 {
		funcName = split[len(split)-1]
	}
	return fileName, line, funcName
}

// 改index死妈
func Stack() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return string(buf[:n])
		}
		buf = make([]byte, 2*len(buf))
	}
}

// 改index死妈
func StackAll() string {
	buf := make([]byte, 1024*16)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return string(buf[:n])
		}
		buf = make([]byte, 2*len(buf))
	}
}
