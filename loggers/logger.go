package loggers

import (
	"fmt"
	"github.com/BlackofZero/go-models/utils"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	LogFolder = utils.GetEnvWithDefault("LOG_FOLDER", "logs")
	LogFile   = utils.GetEnvWithDefault("LOG_File", "opreation")
)

func init() {
	_, err := os.Stat(LogFolder)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(LogFolder, 0755)
	}
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
		go func() {
			for {
				os.Remove(fmt.Sprintf("%s/%s-%s.log", LogFolder, filename, utils.GetDaysAgoStr(7)))
				time.Sleep(11 * time.Hour)
			}
		}()
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
		go func() {
			for {
				os.Remove(fmt.Sprintf("%s/%s-%s.log", LogFolder, filename, utils.GetDaysAgoStr(7)))
				time.Sleep(11 * time.Hour)
			}
		}()
	} else {
		SetLoggerOutPut(logger, fmt.Sprintf("%s/%s.log", LogFolder, filename))
	}
	return logger
}

type BLoggers struct {
	logInfo  *logrus.Logger
	logError *logrus.Logger
	logWarn  *logrus.Logger
}

func Init(fileName, logmode string) *BLoggers {
	bloger := &BLoggers{}
	if logmode == "json" {
		bloger.logInfo = NewJSONFormatterLogger(logrus.InfoLevel, fileName, true)
		bloger.logWarn = NewJSONFormatterLogger(logrus.WarnLevel, fileName, true)
		bloger.logError = NewJSONFormatterLogger(logrus.ErrorLevel, fileName, true)
	} else {
		bloger.logInfo = NewTextFormatterLogger(logrus.InfoLevel, fileName, true)
		bloger.logWarn = NewTextFormatterLogger(logrus.WarnLevel, fileName, true)
		bloger.logError = NewTextFormatterLogger(logrus.ErrorLevel, fileName, true)
	}
	return bloger
}

func (lo *BLoggers) Info(msg string) {
	lo.logInfo.WithFields(logrus.Fields{
		"latency": time.Now(),
	}).Info(msg)

}

func (lo *BLoggers) Error(msg string) {
	lo.logError.WithFields(logrus.Fields{
		"latency": time.Now(),
	}).Error(msg)
}

func (lo *BLoggers) Warnning(msg string) {
	lo.logWarn.WithFields(logrus.Fields{
		"latency": time.Now(),
	}).Warn(msg)
}
