# Marconi Logger

### Overview
Marconi Logger is an application-wide logger and is expected to handle all the logging needs of an application including:
- logging to multiple files
- logging to console
- filtering based on log level

Marconi Logger contains a list of individual log instances, where each one can have its own configurations. Underneath the hood, each log instance is a wrapper around the popular library "logrus".

### Example Usage
```
// initialize logger
func initLogger() {
  mlog.Init("SOME_PATH", "info")
  
  // write to default log
  mlog.GetLogger().Debug("Hello World")
  mlog.GetLogger().Info("Hello World")
  
  // write to file
  logger, _ = mlog.GetLogInstance("FILE_NAME")
  logger.Debug("Hello World")
  logger.Info("Hello World")
}
```

### How to Test
go test

### More info on Logrus
https://github.com/sirupsen/logrus

### Future Work
More improvements are expected:
- log rotation
- setting logging level at file level