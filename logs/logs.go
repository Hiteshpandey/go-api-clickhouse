package logs

import (
	"os"

	"github.com/sirupsen/logrus"
)

type ApiLogger struct {
	logger  *logrus.Logger
	logFile *os.File
}

func New() *ApiLogger {
	lg := &ApiLogger{logger: logrus.New()}
	lg.logger.SetFormatter(&logrus.JSONFormatter{})
	lg.logger.Out = os.Stdout

	file, err := os.OpenFile("velocityApi.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		lg.logFile = file
		lg.logger.Out = file
	} else {
		lg.logger.Error("Failed to log to file, using default stderr", err.Error())
	}
	// }

	return lg
}

func (lg ApiLogger) Write(msg string) {
	lg.logger.Info(msg)
	lg.logFile.WriteString("INFO: " + msg)
}

func (lg ApiLogger) Error(msg, errorMsg string) {
	lg.logger.WithFields(logrus.Fields{
		"Error": errorMsg,
	}).Error(msg)
	lg.logFile.WriteString("ERROR: " + msg + ". Message: " + errorMsg)
}

func (lg ApiLogger) Close(code int) {
	lg.logger.Exit(code)
	lg.logFile.Close()
}
