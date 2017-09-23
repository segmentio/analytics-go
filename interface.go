package analytics

// SegmentClient interface implements Segment functionality
// it extremely useful for testing purposes.
type SegmentClient interface {
	Alias(msg *Alias) error
	Page(msg *Page) error
	Group(msg *Group) error
	Identify(msg *Identify) error
	Track(msg *Track) error
	Close() error
}

// Logger interface imlements logging functions we need.
type Logger interface {
	Printf(format string, v ...interface{})
}
