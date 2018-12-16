package utils

import (
	"path/filepath"
	"runtime"

	"github.com/Sirupsen/logrus"
)

//Locate returns the location of the line calling this function
func Locate() logrus.Fields {
	fields := logrus.Fields{}
	_, path, line, ok := runtime.Caller(1)
	if ok {
		_, file := filepath.Split(path)
		fields["file"] = file
		fields["line"] = line
	}
	return fields
}
