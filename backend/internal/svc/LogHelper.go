package svc

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogHelper struct {
	logger logx.Logger
	ctx    context.Context
}

func NewLogHelper(ctx context.Context) *LogHelper {
	return &LogHelper{
		logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

// Info 记录信息日志
func (l *LogHelper) Info(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Error 记录错误日志
func (l *LogHelper) Error(err error, format string, args ...interface{}) {
	if err != nil {
		format = fmt.Sprintf("%s, error: %v", format, err)
	}
	l.logger.Errorf(format, args...)
}

// Debug 记录调试日志
func (l *LogHelper) Debug(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Slow 记录慢操作日志
func (l *LogHelper) Slow(format string, args ...interface{}) {
	l.logger.Slowf(format, args...)
}
