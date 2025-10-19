package device

type Application interface {
	Init(scheduler func(key string, timeMS int)) error
	TriggerPortInterrupt(port string) error
	TriggerTimeInterrupt(key string) error
}
