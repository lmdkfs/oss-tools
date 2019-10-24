package utils

import (
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func NewLogger() *logrus.Logger {
	if Logger == nil {
		Logger = logrus.New()
	}
	return Logger

}
