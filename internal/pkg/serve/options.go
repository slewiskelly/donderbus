package serve

// Option is an option to a server.
type Option interface {
	apply(*options)
}

// Host sets the host that the server is to serve from.
//
// The default host is 0.0.0.0.
func Host(h string) Option {
	return option(func(o *options) {
		o.host = h
	})
}

// Insecure skips webhook validation.
//
// The default is to perform webhook validation.
func Insecure(ok bool) Option {
	return option(func(o *options) {
		o.insecure = ok
	})
}

// Port sets the port that the server is to listen on.
//
// The default port is 8080.
func Port(p int) Option {
	return option(func(o *options) {
		o.port = p
	})
}

// WebhookSecret sets the secret used to validate webhook requests.
//
// The default secret is the value of GITHUB_WEBHOOK_SECRET.
func WebhookSecret(s []byte) Option {
	return option(func(o *options) {
		o.webhookSecret = s
	})
}

type option func(*options)

func (o option) apply(opts *options) {
	o(opts)
}

type options struct {
	host          string
	insecure      bool
	port          int
	webhookSecret []byte
}
