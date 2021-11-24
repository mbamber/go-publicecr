package publicecr

const (
	defaultBaseURL = "https://gallery.ecr.aws/"
)

// DefaultECRReaderOpts is the default options for an ECRReader
var DefaultECRReaderOpts = &ECRReaderOpts{
	baseURL: defaultBaseURL,
	logger:  nil,
}

// ECRReaderOpts is an option set that can be provided to an ECRReader
type ECRReaderOpts struct {
	baseURL string
	logger  Logger
}

// WithBaseURL returns a functional optional for setting the ECR base URL
func WithBaseURL(url string) func(*ECRReaderOpts) {
	return func(o *ECRReaderOpts) {
		o.baseURL = url
	}
}

// WithDebugLogging returns a functional optional for enabling debug logging
func WithDebugLogger(logger Logger) func(*ECRReaderOpts) {
	return func(o *ECRReaderOpts) {
		o.logger = logger
	}
}
