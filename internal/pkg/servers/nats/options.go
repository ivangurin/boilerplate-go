package nats

type Option func(*server)

func WithName(name string) Option {
	return func(s *server) {
		s.name = name
	}
}

func WithJetStream(domain string) Option {
	return func(s *server) {
		s.withJetStream = true
		s.jetStreamDomain = domain
	}
}

func WithDebug(debug bool) Option {
	return func(s *server) {
		s.debug = debug
	}
}
