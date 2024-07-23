package gotel

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

type TraceLogHandler struct {
	span  trace.Span
	level *slog.Level
	goas  []groupOrAttrs

	next slog.Handler
}

// NewTraceLogHandler creates an instance of [TraceLogHandler] with an optional
// [slog.Logger] and/or tracing specific [slog.LogLevel]
func NewTraceLogHandler(ctx context.Context, span trace.Span, next slog.Handler, level *slog.Level) slog.Handler {
	return &TraceLogHandler{span: span, next: next, level: level}
}

// NewLogWithTracing creates a new [slog.Logger] by adding a TraceLogHandler to
// an existing [slog.Logger]
func NewLogWithTracing(ctx context.Context, span trace.Span, log *slog.Logger) *slog.Logger {
	return slog.New(NewTraceLogHandler(ctx, span, log.Handler(), nil))
}

// NewLogWithTracingAndLevel creates a new [slog.Logger] by adding a TraceLogHandler to
// an existing [slog.Logger] with a custom [slog.Level] specific to tracing
func NewLogWithTracingAndLevel(ctx context.Context, span trace.Span, log *slog.Logger, level slog.Level) *slog.Logger {
	return slog.New(NewTraceLogHandler(ctx, span, log.Handler(), &level))
}

// NewLogForTracing creates a new [slog.Logger] with a TraceLogHandler and a
// specified [slog.Level]
func NewLogForTracing(ctx context.Context, span trace.Span, level slog.Level) *slog.Logger {
	return slog.New(NewTraceLogHandler(ctx, span, nil, &level))
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
// It is called early, before any arguments are processed,
// to save effort if the log event should be discarded.
// If called from a Logger method, the first argument is the context
// passed to that method, or context.Background() if nil was passed
// or the method does not take a context.
// The context is passed so Enabled can use its values
// to make a decision.
//
// For our purposes we are interested in:
// isEnabledTracing(level) || ((next != nil) && next.Enabled(level))
func (h *TraceLogHandler) Enabled(ctx context.Context, level slog.Level) (ok bool) {
	// if we don't want it, does the next handler
	return h.isEnabledTracing(ctx, level) || ((h.next != nil) && h.next.Enabled(ctx, level))
}

// isEnabledTrace determines if the specified level is enabled for Tracing
//
// note that this function is NOT equivalent to IsEnabled because we are only
// interested in the tracing aspect and not the logging aspect.
// The logging aspect is used if the handler does not have an explicit
// level specified. if neither level or next is specified then, [slog.LevelInfo]
// is presumed
func (h *TraceLogHandler) isEnabledTracing(ctx context.Context, level slog.Level) (ok bool) {
	if h.level != nil {
		ok = level >= *h.level
	} else {
		ok = level >= slog.LevelInfo
	}

	// is the span / current span recording?
	if ok {
		span := h.span
		if span == nil {
			// this will always return a span
			span = trace.SpanFromContext(ctx)
		}
		ok = span.IsRecording()
	}

	return ok
}

// Handle handles the Record.
// It will only be called when Enabled returns true.
// The Context argument is as for Enabled.
// It is present solely to provide Handlers access to the context's values.
// Canceling the context should not affect record processing.
// (Among other things, log messages may be necessary to debug a
// cancellation-related problem.)
//
// Handle methods that produce output should observe the following rules:
//
//   - If r.Time is the zero time, ignore the time.
//
//   - If r.PC is zero, ignore it.
//
//   - Attr's values should be resolved.
//
//   - If an Attr's key and value are both the zero value, ignore the Attr.
//     This can be tested with attr.Equal(Attr{}).
//
//   - If a group's key is empty, inline the group's Attrs.
//
//   - If a group has no Attrs (even if it has a non-empty key),
//     ignore it.
//
//     This implementation is co-opted from an example of a simple handler, and
//     performance implications, optimizations, and more sophisticated
//     implementations are noted here:
//     https://github.com/golang/example/blob/master/slog-handler-guide/README.md
//
//     It is likely that we will improve this implementation, but what you see
//     is what you get for now
func (h *TraceLogHandler) Handle(ctx context.Context, r slog.Record) error {
	// if the span is being recorded and the level is appropriate for tracing
	// then log an event / error via the span
	if h.isEnabledTracing(ctx, r.Level) {
		var attrs []attribute.KeyValue

		prefix := ""
		for _, goa := range h.goas {
			if len(goa.group) > 0 {
				prefix = makeKey(prefix, goa.group)
			} else {
				for _, a := range goa.attrs {
					attrs = append(attrs, MakeKeyValue(makeKey(prefix, a.Key), a.Value.Any()))
				}
			}
		}

		if r.NumAttrs() > 0 {
			r.Attrs(func(a slog.Attr) bool {
				attrs = append(attrs, MakeKeyValue(makeKey(prefix, a.Key), a.Value.Any()))
				return true
			})
		}

		// use the associated span, or the span in-effect for the context
		span := h.span
		if span == nil {
			// this will always return a span
			span = trace.SpanFromContext(ctx)
		}

		// use span (not h.span) because it is always non-nil
		if r.Level == slog.LevelError {
			span.RecordError(fmt.Errorf(r.Message), trace.WithAttributes(attrs...))
		} else {
			span.AddEvent(r.Message, trace.WithAttributes(attrs...))
		}
	}

	var err error
	if h.next != nil {
		if h.next.Enabled(ctx, r.Level) {
			err = h.next.Handle(ctx, r)
		}
	}

	return err
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
// The Handler owns the slice: it may retain, modify or discard it.
func (h *TraceLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	return h.withGroupOrAttrs(groupOrAttrs{attrs: attrs})
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
// The keys of all subsequent attributes, whether added by With or in a
// Record, should be qualified by the sequence of group names.
//
// How this qualification happens is up to the Handler, so long as
// this Handler's attribute keys differ from those of another Handler
// with a different sequence of group names.
//
// A Handler should treat WithGroup as starting a Group of Attrs that ends
// at the end of the log event. That is,
//
//	log.WithGroup("s").LogAttrs(ctx, level, msg, slog.Int("a", 1), slog.Int("b", 2))
//
// should behave like
//
//	log.LogAttrs(ctx, level, msg, slog.Group("s", slog.Int("a", 1), slog.Int("b", 2)))
//
// If the name is empty, WithGroup returns the receiver.
func (h *TraceLogHandler) WithGroup(name string) slog.Handler {
	if len(name) == 0 {
		return h
	}

	return h.withGroupOrAttrs(groupOrAttrs{group: name})
}

// groupOrAttrs holds either a group name or a list of slog.Attrs.
type groupOrAttrs struct {
	group string      // group name if non-empty
	attrs []slog.Attr // attrs if non-empty
}

func (h *TraceLogHandler) withGroupOrAttrs(goa groupOrAttrs) *TraceLogHandler {
	h2 := *h
	h2.goas = make([]groupOrAttrs, len(h.goas)+1)
	copy(h2.goas, h.goas)
	h2.goas[len(h2.goas)-1] = goa

	if h.next != nil {
		if len(goa.group) > 0 {
			h2.next = h.next.WithGroup(goa.group)
		} else {
			h2.next = h.next.WithAttrs(goa.attrs)
		}
	}

	return &h2
}

func makeKey(prefix, key string) (result string) {
	if len(prefix) > 0 {
		result = fmt.Sprintf("%s.%s", prefix, key)
	} else {
		result = key
	}

	return
}
