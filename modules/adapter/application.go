package adapter

/*
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include "./include/plugin.h"

// forward declaration for trampoline
extern int   goRead(void *ctx, char* port, char* buf);
extern int   goWrite(void *ctx, char* port, char* buf, int size);
extern void  attachPinInterrupt(void *ctx, int pin, interrupt_callback_t cb);
extern void  attachTimeInterrupt(void *ctx, int time_ms, short periodic, interrupt_callback_t cb);
extern void* dataGetter(void *ctx, char* name);
extern void  dataSetter(void *ctx, char* name, void* value);
extern void  goLog(char *line);

// trampoline wrapper
static void tLibInit(lib_init_func_t lib_init) {
	interface_t iface = {
		.read_port             = goRead,
		.write_port            = goWrite,
		.attach_pin_interrupt  = attachPinInterrupt,
		.attach_time_interrupt = attachTimeInterrupt,
		.data_getter           = dataGetter,
		.data_setter           = dataSetter,
		.log                   = goLog,
	};
    lib_init(iface);
}

static void tInit(void *ctx, init_func_t init) {
    init(ctx);
}
static void tShutdown(void *ctx, shutdown_t shutdown) {
	shutdown(ctx);
}
static void tInterrupt(void *ctx, interrupt_callback_t cb) {
	cb(ctx);
}
*/
import "C"

import (
	"fmt"
	"runtime/cgo"
	"unsafe"

	"github.com/Gordy96/cgo_dl/dl"
	"github.com/Gordy96/evt-sim/simulation"
)

//export goLog
func goLog(line *C.char) {
	fmt.Printf("%s\n", C.GoString(line))
}

type pinInterruptConfig struct {
	pin int
	cb  C.interrupt_callback_t
}

//export attachPinInterrupt
func attachPinInterrupt(ctx *C.void, pin C.int, cb C.interrupt_callback_t) {
	c := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	cfg := pinInterruptConfig{
		pin: int(pin),
		cb:  cb,
	}
	c.pinInterrupts[fmt.Sprintf("%d", int(pin))] = cfg
}

type timerInterruptConfig struct {
	timeMS   int
	periodic bool
	cb       C.interrupt_callback_t
}

//export attachTimeInterrupt
func attachTimeInterrupt(ctx *C.void, timeMS C.int, periodic C.short, cb C.interrupt_callback_t) {
	c := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	cfg := timerInterruptConfig{
		timeMS:   int(timeMS),
		periodic: int(periodic) > 0,
		cb:       cb,
	}
	key := fmt.Sprintf("%d", int(timeMS))
	if cfg.periodic {
		key += "_periodic"
	}
	c.timerInterrupts[key] = cfg
	c.schedule(key, cfg.timeMS)
}

//export goRead
func goRead(ctx *C.void, port *C.char, buf *C.char) C.int {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)

	if p, ok := a.ports[C.GoString(port)]; ok {
		temp := make([]byte, 1024)
		n, _ := p.Read(temp)

		C.memcpy(unsafe.Pointer(buf), unsafe.Pointer(&temp[0]), C.size_t(n))

		return C.int(n)
	}

	return 0
}

//export goWrite
func goWrite(ctx *C.void, port *C.char, buf *C.char, size C.int) C.int {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)

	if p, ok := a.ports[C.GoString(port)]; ok {
		temp := make([]byte, int(size))
		C.memcpy(unsafe.Pointer(&temp[0]), unsafe.Pointer(buf), C.size_t(size))
		n, _ := p.Write(temp)

		return C.int(n)
	}

	return 0
}

//export dataGetter
func dataGetter(ctx *C.void, name *C.char) *C.void {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	v, _ := a.mem[C.GoString(name)].(unsafe.Pointer)
	return (*C.void)(v)
}

//export dataSetter
func dataSetter(ctx *C.void, name *C.char, value *C.void) {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	a.mem[C.GoString(name)] = unsafe.Pointer(value)
}

type Port interface {
	Name() string
	Write([]byte) (int, error)
	Read([]byte) (int, error)
}

type Application struct {
	simulation.ParameterBag
	selfUnsafe      cgo.Handle
	id              string
	initFunc        C.init_func_t
	shutdownFunc    C.shutdown_t
	pinInterrupts   map[string]pinInterruptConfig
	timerInterrupts map[string]timerInterruptConfig
	ports           map[string]Port
	mem             map[string]interface{}
	schedule        func(string, int)
}

func (a *Application) Close() error {
	C.tShutdown(unsafe.Pointer(a.selfUnsafe), a.shutdownFunc)
	a.selfUnsafe.Delete()
	return nil
}

func (a *Application) ID() string {
	return a.id
}

func (a *Application) Init() {
	C.tInit(unsafe.Pointer(a.selfUnsafe), a.initFunc)
}

func (a *Application) TriggerTimeInterrupt(key string) {
	if i, ok := a.timerInterrupts[key]; ok {
		C.tInterrupt(unsafe.Pointer(a.selfUnsafe), i.cb)
		if i.periodic {
			a.schedule(key, i.timeMS)
		}
	}
}

func (a *Application) TriggerPinInterrupt(pin int) error {
	key := fmt.Sprintf("%d", pin)
	if i, ok := a.pinInterrupts[key]; ok {
		C.tInterrupt(unsafe.Pointer(a.selfUnsafe), i.cb)
	}

	return nil
}

func New(id string, ports []Port, schedule func(string, int), lib *dl.SO) (*Application, error) {
	sym, err := lib.Func("init")
	if err != nil {
		return nil, err
	}

	initFunc := (C.init_func_t)(sym)

	sym, err = lib.Func("shutdown")
	if err != nil {
		return nil, err
	}

	shutdownFunc := (C.shutdown_t)(sym)

	a := &Application{
		id:              id,
		initFunc:        initFunc,
		shutdownFunc:    shutdownFunc,
		timerInterrupts: make(map[string]timerInterruptConfig),
		pinInterrupts:   make(map[string]pinInterruptConfig),
		ports:           make(map[string]Port),
		mem:             make(map[string]interface{}),
		schedule:        schedule,
	}

	a.selfUnsafe = cgo.NewHandle(a)

	for _, port := range ports {
		a.ports[port.Name()] = port
	}

	return a, nil
}

func OpenLib(path string) (*dl.SO, error) {
	lib, err := dl.Open(path)
	if err != nil {
		return nil, err
	}

	sym, err := lib.Func("init_lib")
	if err != nil {
		return nil, err
	}

	initFunc := (C.lib_init_func_t)(sym)
	C.tLibInit(initFunc)

	return lib, nil
}
