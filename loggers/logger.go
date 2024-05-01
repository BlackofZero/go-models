package loggers

import (
	"fmt"
	"github.com/BlackofZero/go-models/utils"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	LogFolder     = utils.GetEnvWithDefault("LOG_FOLDER", "logs")
	LogFile       = utils.GetEnvWithDefault("LOG_File", "opreation")
	Logger        BLoggers
	RequestLogger BLoggers
)

func init() {
	_, err := os.Stat(LogFolder)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(LogFolder, 0755)
	}
	Logger = BLoggers{fileName: LogFile}
	RequestLogger = BLoggers{fileName: "request"}
}

func SetLoggerOutPut(logger *logrus.Logger, filename string) {
	file, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	logger.Out = file
	//if configs.DEBUG {
	//	logger.SetOutput(io.MultiWriter(os.Stdout, file))
	//} else {
	//	logger.Out = file
	//}
}

func NewJSONFormatterLogger(level logrus.Level, filename string, rotate bool) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(level)
	if rotate {
		go func() {
			for {
				SetLoggerOutPut(logger, fmt.Sprintf("%s/%s-%s.log", LogFolder, filename, utils.GetNowDayStr()))
				time.Sleep(time.Minute)
			}
		}()
		//go func() {
		//	for {
		//		os.Remove(fmt.Sprintf("%s/%s-%s.log", LogFolder, filename, utils.GetDaysAgoStr(7)))
		//		time.Sleep(11 * time.Hour)
		//	}
		//}()
	} else {
		SetLoggerOutPut(logger, fmt.Sprintf("%s/%s.log", LogFolder, filename))
	}
	return logger
}

func NewTextFormatterLogger(level logrus.Level, filename string, rotate bool) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetLevel(level)
	if rotate {
		go func() {
			for {
				SetLoggerOutPut(logger, fmt.Sprintf("%s/%s-%s.log", LogFolder, filename, utils.GetNowDayStr()))
				time.Sleep(time.Minute)
			}
		}()
		//go func() {
		//	for {
		//		os.Remove(fmt.Sprintf("%s/%s-%s.log", LogFolder, filename, utils.GetDaysAgoStr(7)))
		//		time.Sleep(11 * time.Hour)
		//	}
		//}()
	} else {
		SetLoggerOutPut(logger, fmt.Sprintf("%s/%s.log", LogFolder, filename))
	}
	return logger
}

type BLoggers struct {
	fileName string
}

func (lo BLoggers) Info(msg string) {
	Log := NewJSONFormatterLogger(logrus.InfoLevel, lo.fileName, true)
	Log.WithFields(logrus.Fields{
		"latency": time.Now(),
	}).Info(msg)

}

func (lo BLoggers) Error(msg string) {
	Log := NewJSONFormatterLogger(logrus.ErrorLevel, lo.fileName, true)
	Log.WithFields(logrus.Fields{
		"latency": time.Now(),
	}).Error(msg)
}

func (lo BLoggers) Warnning(msg string) {
	Log := NewJSONFormatterLogger(logrus.InfoLevel, lo.fileName, true)
	Log.WithFields(logrus.Fields{
		"latency": time.Now(),
	}).Warn(msg)
}
func (lo BLoggers) Write(p []byte) (n int, err error) {
	return utils.WriteBytes("logs/box_image", p)
}
