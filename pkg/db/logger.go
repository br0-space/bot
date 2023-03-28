package db

import (
	"context"
	"errors"
	"fmt"
	logger "github.com/br0-space/bot-logger"
	"github.com/br0-space/bot/interfaces"
	gormLogger "gorm.io/gorm/logger"
	"regexp"
	"runtime"
	"time"
)

type gormLoggerBridge struct {
	wrappedLogger interfaces.LoggerInterface
	config        gormLogger.Config
}

func NewGormLoggerBridge(wrappedLogger logger.Interface) gormLogger.Interface {
	return &gormLoggerBridge{
		wrappedLogger: wrappedLogger,
		config: gormLogger.Config{
			SlowThreshold:             10000 * time.Millisecond,
			LogLevel:                  gormLogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	}
}

func (l *gormLoggerBridge) LogMode(_ gormLogger.LogLevel) gormLogger.Interface {
	// Ignore (log level is set elsewhere)
	return l
}

func (l gormLoggerBridge) Info(_ context.Context, msg string, data ...interface{}) {
	l.wrappedLogger.SetExtraCallDepth(l.getExtraCallDepth())
	l.wrappedLogger.Infof(msg, data...)
	l.wrappedLogger.ResetExtraCallDepth()
}

func (l gormLoggerBridge) Warn(_ context.Context, msg string, data ...interface{}) {
	l.wrappedLogger.SetExtraCallDepth(l.getExtraCallDepth())
	l.wrappedLogger.Warningf(msg, data...)
	l.wrappedLogger.ResetExtraCallDepth()
}

func (l gormLoggerBridge) Error(_ context.Context, msg string, data ...interface{}) {
	l.wrappedLogger.SetExtraCallDepth(l.getExtraCallDepth())
	l.wrappedLogger.Errorf(msg, data...)
	l.wrappedLogger.ResetExtraCallDepth()
}

func (l gormLoggerBridge) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	l.wrappedLogger.SetExtraCallDepth(l.getExtraCallDepth())

	// Stolen from https://github.com/op/go-logging/blob/master/logger.go
	elapsed := time.Since(begin)
	switch {
	case err != nil && (!errors.Is(err, gormLogger.ErrRecordNotFound) || !l.config.IgnoreRecordNotFoundError):
		format := "%s [%.3fms] [rows:%v] %s"
		sql, rows := fc()
		if rows == -1 {
			l.wrappedLogger.Errorf(format, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.wrappedLogger.Errorf(format, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.config.SlowThreshold && l.config.SlowThreshold != 0:
		format := "%s\n[%.3fms] [rows:%v] %s"
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.config.SlowThreshold)
		if rows == -1 {
			l.wrappedLogger.Warningf(format, slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.wrappedLogger.Warningf(format, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	default:
		format := "[%.3fms] [rows:%v] %s"
		sql, rows := fc()
		if rows == -1 {
			l.wrappedLogger.Debugf(format, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.wrappedLogger.Debugf(format, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}

	l.wrappedLogger.ResetExtraCallDepth()
}

func (l gormLoggerBridge) getExtraCallDepth() int {
	extraCallDepth := 1
	// Stolen from https://github.com/go-gorm/gorm/blob/master/utils/utils.go
	re := regexp.MustCompile(`gorm.io/gorm`)
	for i := 2; i < 15; i++ {
		_, file, _, _ := runtime.Caller(i)
		if match := re.MatchString(file); match {
			extraCallDepth++
		} else {
			break
		}
	}
	return extraCallDepth
}
