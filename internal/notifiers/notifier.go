package notifiers

type Notifier interface {
	Notify(message string) error
}
