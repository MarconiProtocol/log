package mlog

import (
  "fmt"
  "github.com/sirupsen/logrus"
  "io"
  "os"
  "path/filepath"
  "sync"
)

// Marconi Log
type Mlog struct {
  name     string
  file     *os.File
  logger   *logrus.Logger
}

// Struct containing all log instances
type mlogs struct {
  logMap map[string]*Mlog
  mutex  sync.Mutex
}

var mLogs *mlogs
var defaultLogger *Mlog
var mLogDir = "."
var outputLevel = "info"

func (mlog *Mlog) Debug(msg ...interface{}) {
  mlog.logger.Debug(msg...)
}

func (mlog *Mlog) Debugf(format string, v ...interface{}) {
  mlog.logger.Debugf(format, v...)
}

func (mlog *Mlog) Info(msg ...interface{}) {
  mlog.logger.Info(msg...)
}

func (mlog *Mlog) Infof(format string, v ...interface{}) {
  mlog.logger.Infof(format, v...)
}

func (mlog *Mlog) Warn(msg ...interface{}) {
  mlog.logger.Warn(msg...)
}

func (mlog *Mlog) Warnf(format string, v ...interface{}) {
  mlog.logger.Warnf(format, v...)
}

func (mlog *Mlog) Warning(msg ...interface{}) {
  mlog.logger.Warn(msg...)
}

func (mlog *Mlog) Warningf(format string, v ...interface{}) {
  mlog.logger.Warnf(format, v...)
}

func (mlog *Mlog) Error(msg ...interface{}) {
  mlog.logger.Error(msg...)
}

func (mlog *Mlog) Errorf(format string, v ...interface{}) {
  mlog.logger.Errorf(format, v...)
}

func (mlog *Mlog) Fatal(msg ...interface{}) {
  mlog.logger.Fatal(msg...)
}

func (mlog *Mlog) Fatalf(format string, v ...interface{}) {
  mlog.logger.Fatalf(format, v...)
}

// initialize the logging directory and log level for all log instances
func Init(logDir string, level string) {
  fmt.Println("init logDir", logDir)
  // keep a map of loggers for file specific logging
  mLogs = new(mlogs)
  mLogs.logMap = make(map[string]*Mlog)
  if logDir != "" {
    mLogDir = logDir
  }
  if level != "" {
    outputLevel = level
  }

  // init default logger for all non-file specific logs
  var err error
  defaultLogger, err = GetLogInstance("application", true)
  if err != nil {
    fmt.Println("failed to init: log - application", err)
    os.Exit(1)
  }
}

// Get the log instance by name, if not exist, create a new one, optional parameter specifies
// whether the log instance outputs to stdout as well as to file or just to file
func GetLogInstance(name string, optional ...bool) (mlog *Mlog, err error) {
  if mLogs == nil {
    Init(mLogDir, outputLevel)
  }
  mLogs.mutex.Lock()
  defer mLogs.mutex.Unlock()

  m, found := mLogs.logMap[name]
  if !found {
    // init a new logger and add it to map
    mlog = new(Mlog)
    mlog.name = name

    stdout := false
    if len(optional) > 0 {
      stdout = optional[0]
    }
    err = mlog.open(stdout)
    mLogs.logMap[name] = mlog
  } else {
    mlog = m
  }
  return mlog, err
}

// Returns the default logger
func GetLogger() (mlog *Mlog) {
  if defaultLogger == nil {
    Init(mLogDir, outputLevel)
  }
  return defaultLogger
}

// Open the file of a log instance and prepare it for writing
func (mlog *Mlog) open(stdout bool) (err error) {
  // Create directories if they don't already exist
  err = os.MkdirAll(mLogDir, 0777)
  if err != nil {
    fmt.Println("failed to open dir: ", mLogDir)
    return err
  }

  filePath := filepath.Join(mLogDir, mlog.name + ".log")
  file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
  if err != nil {
    fmt.Println("failed to open: ", filePath)
    return err
  }
  mlog.file = file
  mlog.logger = logrus.New()
  if stdout {
    multiWriter := io.MultiWriter(file, os.Stdout)
    mlog.logger.SetOutput(multiWriter)
  } else {
    mlog.logger.SetOutput(file)
  }
  logLevel, err := logrus.ParseLevel(outputLevel)
  if err != nil {
    fmt.Println("failed to parse log level")
    return err
  }
  mlog.logger.SetLevel(logLevel)
  return nil
}

// Set the log level of all log instances
func SetOutputLevel(level string) {
  outputLevel = level

  mLogs.mutex.Lock()
  defer mLogs.mutex.Unlock()

  logLevel, err := logrus.ParseLevel(outputLevel)
  if err == nil {
    defaultLogger.logger.SetLevel(logLevel)
    for _, v := range mLogs.logMap {
      v.logger.SetLevel(logLevel)
    }
  }
}

// Close the file of a log instance from writing
func (mlog *Mlog) Close() {
  mlog.file.Close()
}
