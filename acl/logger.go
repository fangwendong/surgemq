package acl

import "go.uber.org/zap"

var logger *zap.Logger

func init() {
	logger, _ = zap.NewDevelopment()
}

func SetLogger(l *zap.Logger) {
	logger = l
	logger.Named("surgemq")
}
