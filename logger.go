package jasmine

import (
	"github.com/charmbracelet/log"
	"github.com/rs/zerolog"
	log2 "github.com/rs/zerolog/log"
	"os"
	"strings"
)

var Logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller:    true,
	ReportTimestamp: true,
	Prefix:          "Jasmine",
})

type zeroToCharmLogger struct{}

func (z *zeroToCharmLogger) Write(p []byte) (n int, err error) {
	s := string(p)
	if len(s) <= 3 {
		return len(p), nil
	}
	inLevel := s[:3]
	var level log.Level
	switch inLevel {
	case "TRC", "DBG":
		level = log.DebugLevel
	case "INF":
		level = log.InfoLevel
	case "WRN":
		level = log.WarnLevel
	case "ERR":
		level = log.ErrorLevel
	case "FTL", "PNC":
		level = log.FatalLevel
	}
	Logger.Helper()
	translator := Logger.With()
	translator.SetCallerFormatter(func(s string, i int, s2 string) string {
		return "tokidoki internal"
	})
	translator.Log(level, strings.TrimSpace(s[4:]))
	return len(p), nil
}

func init() {
	log2.Logger = log2.Output(zerolog.ConsoleWriter{
		Out: &zeroToCharmLogger{},
		FormatTimestamp: func(i interface{}) string {
			return ""
		},
	})
}
