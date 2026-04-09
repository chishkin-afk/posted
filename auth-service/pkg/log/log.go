package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"log/slog"
)

type Color string

const (
	ColorReset  Color = "\033[0m"
	ColorRed    Color = "\033[31m"
	ColorGreen  Color = "\033[32m"
	ColorYellow Color = "\033[33m"
	ColorBlue   Color = "\033[34m"
	ColorPurple Color = "\033[35m"
	ColorCyan   Color = "\033[36m"
	ColorGray   Color = "\033[37m"
	ColorWhite  Color = "\033[97m"
)

const (
	EnvDev   = "dev"
	EnvLocal = "local"
	EnvProd  = "prod"
)

func New(env string) *slog.Logger {
	var handler slog.Handler
	var level slog.Level

	switch env {
	case EnvDev, EnvLocal:
		level = slog.LevelDebug
		handler = NewColoredHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		})
	case EnvProd:
		level = slog.LevelInfo
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		})
	default:
		level = slog.LevelDebug
		handler = NewColoredHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		})
	}

	return slog.New(handler)
}

type ColoredHandler struct {
	mu     sync.Mutex
	out    io.Writer
	opts   *slog.HandlerOptions
	attrs  []slog.Attr
	groups []string
}

func NewColoredHandler(out io.Writer, opts *slog.HandlerOptions) *ColoredHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &ColoredHandler{
		out:  out,
		opts: opts,
	}
}

func (h *ColoredHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *ColoredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.attrs = append(h.attrs, attrs...)
	return h
}

func (h *ColoredHandler) WithGroup(name string) slog.Handler {
	h.groups = append(h.groups, name)
	return h
}

func (h *ColoredHandler) Handle(ctx context.Context, record slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	timeStr := record.Time.Format("2006-01-02 15:04:05")

	var levelColor Color
	var levelStr string
	switch record.Level {
	case slog.LevelDebug:
		levelColor = ColorGray
		levelStr = "DEBUG"
	case slog.LevelInfo:
		levelColor = ColorGreen
		levelStr = "INFO"
	case slog.LevelWarn:
		levelColor = ColorYellow
		levelStr = "WARN"
	case slog.LevelError:
		levelColor = ColorRed
		levelStr = "ERROR"
	default:
		levelColor = ColorWhite
		levelStr = record.Level.String()
	}

	var sourceStr string
	if record.PC != 0 && h.opts.AddSource {
		fs := runtime.CallersFrames([]uintptr{record.PC})
		frame, _ := fs.Next()
		if frame.File != "" {
			shortFile := frame.File
			idx := strings.LastIndex(shortFile, "/")
			if idx != -1 {
				shortFile = shortFile[idx+1:]
			}
			sourceStr = fmt.Sprintf("%s%s:%d%s", ColorCyan, shortFile, frame.Line, ColorReset)
		}
	}

	var buf strings.Builder
	for _, attr := range h.attrs {
		buf.WriteString(" ")
		buf.WriteString(attr.Key)
		buf.WriteString("=")
		buf.WriteString(fmt.Sprintf("%v", attr.Value.Any()))
	}
	record.Attrs(func(a slog.Attr) bool {
		buf.WriteString(" ")
		buf.WriteString(a.Key)
		buf.WriteString("=")
		buf.WriteString(fmt.Sprintf("%v", a.Value.Any()))
		return true
	})
	extraAttrs := buf.String()

	line := fmt.Sprintf("%s [%s%s%s] %s %s%s%s %s\n",
		timeStr,
		levelColor, levelStr, ColorReset,
		sourceStr,
		ColorWhite, record.Message, ColorReset,
		extraAttrs,
	)

	_, err := h.out.Write([]byte(line))
	return err
}
