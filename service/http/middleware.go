package http

import (
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

// logging
func NewFilteredLogger(excluded []string) middleware.LogFormatter {
	if excluded == nil {
		excluded = make([]string, 0)
	}
	return &filteredLogger{
		excluded: excluded,
		logger:   &middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags)},
	}
}

type noopLogEntry struct{}

func (noopLogEntry) Write(status, bytes int, elapsed time.Duration) {}
func (noopLogEntry) Panic(v interface{}, stack []byte)              {}

type filteredLogger struct {
	excluded []string
	logger   middleware.LogFormatter
}

func (this *filteredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	for _, ex := range this.excluded {
		if ex == r.URL.Path {
			return noopLogEntry{}
		}
	}
	return this.logger.NewLogEntry(r)
}
