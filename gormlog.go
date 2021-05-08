package gorm2logrus

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type GormLogrus struct {
	logrus.Logger
	// for gorm v2
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
}

func NewGormLogger() *GormLogrus {
	return &GormLogrus{
		Logger: logrus.Logger{
			Out:          os.Stderr,
			Formatter:    new(logrus.TextFormatter),
			Hooks:        make(logrus.LevelHooks),
			Level:        logrus.TraceLevel,
			ExitFunc:     os.Exit,
			ReportCaller: false,
		},
		SkipErrRecordNotFound: false,
	}
}

func (l *GormLogrus) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	// l.SetLevel(Level(level))
	return l
}

func (l *GormLogrus) Info(ctx context.Context, s string, args ...interface{}) {
	l.WithContext(ctx)
	l.Log(logrus.InfoLevel, fmt.Sprintf(s, args...))
	// log.WithContext(ctx).Infof(s, args)
	// log "github.com/sirupsen/logrus"
}

func (l *GormLogrus) Warn(ctx context.Context, s string, args ...interface{}) {
	l.WithContext(ctx)
	l.Log(logrus.WarnLevel, fmt.Sprintf(s, args...))
}

func (l *GormLogrus) Error(ctx context.Context, s string, args ...interface{}) {
	l.WithContext(ctx)
	l.Log(logrus.ErrorLevel, fmt.Sprintf(s, args...))
}

func (l *GormLogrus) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrus.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}

	traceLevel := logrus.DebugLevel

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		traceLevel = logrus.ErrorLevel
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		traceLevel = logrus.WarnLevel
	}

	l.WithContext(ctx).WithFields(fields)

	l.Log(traceLevel, fmt.Sprintf("%s [%v]", sql, elapsed))
}

func (l *GormLogrus) SetSkipErrRecordNotFound(SkipErrRecordNotFound bool) gormlogger.Interface {
	l.SkipErrRecordNotFound = SkipErrRecordNotFound
	return l
}

func (l *GormLogrus) SetLogMode(level logrus.Level) gormlogger.Interface {
	l.SetLevel(level)
	return l
}
