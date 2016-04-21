package analytics

type Message interface {
	validate() error

	serializable() interface{}
}
