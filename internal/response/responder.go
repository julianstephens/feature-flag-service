package response

// New creates a new Responder with default JSON encoder and hooks.
// The default configuration includes:
// - JSONEncoder for JSON responses
// - DefaultBefore, DefaultAfter (no-op hooks)
// - DefaultOnError for error handling
func New() *Responder {
	return &Responder{
		Encoder: NewJSONEncoder(),
		Before:  DefaultBefore,
		After:   DefaultAfter,
		OnError: DefaultOnError,
	}
}

// NewWithLogging creates a new Responder with JSON encoder and logging hooks.
// This configuration logs requests, responses, and errors.
func NewWithLogging() *Responder {
	return &Responder{
		Encoder: NewJSONEncoder(),
		Before:  LoggingBefore,
		After:   LoggingAfter,
		OnError: LoggingOnError,
	}
}

// NewCustom creates a new Responder with custom configuration.
// All parameters are optional and will use defaults if nil.
func NewCustom(encoder Encoder, before BeforeFunc, after AfterFunc, onError OnErrorFunc) *Responder {
	if encoder == nil {
		encoder = NewJSONEncoder()
	}
	if before == nil {
		before = DefaultBefore
	}
	if after == nil {
		after = DefaultAfter
	}
	if onError == nil {
		onError = DefaultOnError
	}

	return &Responder{
		Encoder: encoder,
		Before:  before,
		After:   after,
		OnError: onError,
	}
}