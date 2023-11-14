package sl

import (
	"log/slog"

	_ "github.com/lib/pq"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
