package chttp

import "github.com/hashicorp/go-retryablehttp"

type (
	// Logger interface allows to use other loggers than
	// standard log.Logger.
	Logger = retryablehttp.Logger

	// LeveledLogger is an interface that can be implemented by any logger or a
	// logger wrapper to provide leveled logging. The methods accept a message
	// string and a variadic number of key-value pairs. For log.Printf style
	// formatting where message string contains a format specifier, use Logger
	// interface.
	LeveledLogger = retryablehttp.LeveledLogger

	// CheckRetry specifies a policy for handling retries. It is called
	// following each request with the response and error values returned by
	// the http.Client. If CheckRetry returns false, the Client stops retrying
	// and returns the response to the caller. If CheckRetry returns an error,
	// that error value is returned in lieu of the error from the request. The
	// Client will close any response body when retrying, but if the retry is
	// aborted it is up to the CheckRetry callback to properly close any
	// response body before returning.
	CheckRetry = retryablehttp.CheckRetry

	// Backoff specifies a policy for how long to wait between retries.
	// It is called after a failing request to determine the amount of time
	// that should pass before trying again.
	Backoff = retryablehttp.Backoff

	// ErrorHandler is called if retries are expired, containing the last status
	// from the http library. If not specified, default behavior for the library is
	// to close the body and return an error indicating how many tries were
	// attempted. If overriding this, be sure to close the body if needed.
	ErrorHandler = retryablehttp.ErrorHandler

	internalRequest = retryablehttp.Request
)
