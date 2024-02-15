package logger

import (
	"github.com/julioisaac/users/config"
	"github.com/sirupsen/logrus"
)

var Logger = logrus.WithFields(logrus.Fields{
	"app": config.GetString("SERVICE_NAME"),
	"env": config.GetString("ENV"),
})

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyLevel: "severity",
		},
	})

	logrus.RegisterExitHandler(func() {
		logrus.Info("application will stop probably due to a os signal")
	})

	ll := config.GetString("LOG_LEVEL")
	l, err := logrus.ParseLevel(ll)
	if err != nil {
		Logger.WithError(err).Errorf("error parsing log level %s", ll)
		return
	}

	logrus.SetLevel(l)
}
