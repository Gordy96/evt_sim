package device

import "time"

type Port interface {
	Name() string
	Write([]byte) (int, error)
	Read([]byte) (int, error)
}

type Application interface {
	Init(scheduler func(key string, timeMS int), ports ...Port) error
	TriggerPortInterrupt(port string) error
	TriggerTimeInterrupt(key string, now time.Duration) error
	Close() error
}
