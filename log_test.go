package mlog

import (
  "fmt"
  "github.com/stretchr/testify/assert"
  "io/ioutil"
  "os"
  "testing"
)

// test logging via default logger
func TestDefaultLogger(t *testing.T) {
  mlog := GetLogger()
  assert.NotNil(t, mlog)

  // write some dummy logs
  var logExpected = "testing testing"
  mlog.Info(logExpected)

  byteArray, err := ioutil.ReadFile("application.log")
  if err != nil {
    fmt.Print(err)
  }
  logActual := string(byteArray)

  assert.Contains(t, logActual, logExpected, "should be in output")
  assert.NotContains(t, logActual, "random string", "should not be in output")
}

// test logging to a specific file
func TestFileLogger(t *testing.T) {
  Init(".", "info")

  mlog, err := GetLogInstance("test_file")
  assert.Nil(t, err)

  // write some dummy logs
  var logExpected = "test message"
  mlog.Info(logExpected)

  byteArray, err := ioutil.ReadFile("test_file.log")
  if err != nil {
    fmt.Print(err)
  }
  logActual := string(byteArray)

  assert.Contains(t, logActual, logExpected, "should be in output")
  assert.NotContains(t, logActual, "another random string", "should not be in output")
}

// test for a higher log level
func TestLogLevel(t *testing.T) {
  Init(".", "warn")

  mlog, err := GetLogInstance("test_log_level")
  assert.Nil(t, err)

  // write some dummy logs
  var warnMsg = "warn message"
  mlog.Warn(warnMsg)
  var infoMsg = "info message"
  mlog.Info(infoMsg)

  byteArray, err := ioutil.ReadFile("test_log_level.log")
  if err != nil {
    fmt.Print(err)
  }
  logActual := string(byteArray)

  assert.Contains(t, logActual, warnMsg, "should be in output")
  assert.NotContains(t, logActual, infoMsg, "should not be in output")
}

// test for changing of log level after logger is created
func TestLogLevelChanges(t *testing.T) {
  Init(".", "info")

  mlog, err := GetLogInstance("test_log_level_changes")
  assert.Nil(t, err)

  // write some dummy logs
  var msg1 = "message 1"
  var msg2 = "message 2"
  var msg3 = "message 3"
  mlog.Info(msg1)
  mlog.Info(msg2)
  SetOutputLevel("warn")
  mlog.Info(msg3)

  byteArray, err := ioutil.ReadFile("test_log_level_changes.log")
  if err != nil {
    fmt.Print(err)
  }
  logActual := string(byteArray)

  assert.Contains(t, logActual, msg1, "should be in output")
  assert.Contains(t, logActual, msg2, "should be in output")
  assert.NotContains(t, logActual, msg3, "should not be in output")

  // only need this in the final test
  Cleanup()
}

// clean up the log files generated from the testing
func Cleanup() {
  os.Remove("application.log")
  os.Remove("test_file.log")
  os.Remove("test_log_level.log")
  os.Remove("test_log_level_changes.log")
}