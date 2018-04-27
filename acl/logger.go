package acl

import (
	"go.uber.org/zap"
)

var Logger = NewLogger()

func NewLogger() *zap.Logger {
	l, _ := zap.NewProductionConfig().Build()
	return l
}
