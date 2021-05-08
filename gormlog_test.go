package gorm2logrus

import (
	"testing"
)

func TestGormLogger_Info(t *testing.T) {
	logg := NewGormLogger()
	logg.Info(nil, "hello %s", "logrus world")
}
