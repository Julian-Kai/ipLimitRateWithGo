package log

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func Initial() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05.000",
		DisableHTMLEscape: true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "0_level",
			logrus.FieldKeyTime:  "0_time",
			logrus.FieldKeyFile:  "1_file",
			logrus.FieldKeyFunc:  "1_func",
			logrus.FieldKeyMsg:   "1_msg",
		},
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			_, fileName := path.Split(f.File)
			return fmt.Sprintf("%s()", funcName), fmt.Sprintf("%s:%d", fileName, f.Line)
		},
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the info severity or above.
	logrus.SetLevel(logrus.InfoLevel)

	// Show method name
	logrus.SetReportCaller(true)
}
