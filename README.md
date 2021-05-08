# gorm2logrus

**logrus for gorm v2 base on sirupsen/logrus && onrik/gorm-logrus** 

## example
```golang
package main

import github.com/FWangZil/gorm2logrus

// Logger can add hooks and use them normally globally, supporting all the features that Logur itself supports
var Logger *logrus.Logger
// Custom Logger for use by GORM v2
var GormLogger *gorm2logrus.GormLogger

var DB *gorm.DB

func main(){
    initLogrus()
    initDataBase(GormLogger)
}

func initLogrus() {
	logger := gorm2logrus.NewGormLogger()

	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})

	GormLogger = logger
	Logger = &GormLogger.Logger
}

func initDataBase(gormLogger *gorm2logrus.GormLogger) {
    var err error
    
    gormConfig := &gorm.Config{
        DisableForeignKeyConstraintWhenMigrating: true, 
	}
	if gormLogger != nil {
        gormConfig.Logger = gormLogger
	}

	// https://github.com/go-gorm/postgres
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
        // disables implicit prepared statement usage
		PreferSimpleProtocol: true, 
	}), gormConfig)
	if err != nil {
		Logger.Fatalf("pgsql connect failed: %v", err)
	}
}

```