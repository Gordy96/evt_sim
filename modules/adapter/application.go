package adapter

/*
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include "./include/plugin.h"

typedef const char cchar_t;

// forward declaration for trampoline
extern int   goRead(void *ctx, const char* port, char* buf, size_t size);
extern int   goWrite(void *ctx, const char* port, char* buf, size_t size);
extern void  attachPortInterrupt(void *ctx, const char* port, interrupt_callback_t cb);
extern void  attachTimeInterrupt(void *ctx, int time_ms, short periodic, interrupt_callback_t cb);
extern void* dataGetter(void *ctx, const char* name);
extern void  dataSetter(void *ctx, const char* name, void* value);
extern int   stringParamGetter(void *ctx, const char* name, char* buf, size_t size);
extern int   int8ParamGetter(void *ctx, const char* name, int8_t* dst);
extern int   int16ParamGetter(void *ctx, const char* name, int16_t* dst);
extern int   int32ParamGetter(void *ctx, const char* name, int32_t* dst);
extern int   int64ParamGetter(void *ctx, const char* name, int64_t* dst);
extern int   uint8ParamGetter(void *ctx, const char* name, uint8_t* dst);
extern int   uint16ParamGetter(void *ctx, const char* name, uint16_t* dst);
extern int   uint32ParamGetter(void *ctx, const char* name, uint32_t* dst);
extern int   uint64ParamGetter(void *ctx, const char* name, uint64_t* dst);
extern int   doubleParamGetter(void *ctx, const char* name, double* dst);
extern void  goLog(void *ctx, int level, char *line);
extern void  goPacketDump(void *ctx, const char* dir, void *data, size_t size);

// trampoline wrapper
static void tLibInit(lib_init_func_t lib_init) {
	interface_t iface = {
		.read_port             = goRead,
		.write_port            = goWrite,
		.attach_port_interrupt = attachPortInterrupt,
		.attach_time_interrupt = attachTimeInterrupt,
		.get_data              = dataGetter,
		.set_data              = dataSetter,
		.log                   = goLog,
		.dump_packet           = goPacketDump,
        .get_string_param      = stringParamGetter,
        .get_int8_param        = int8ParamGetter,
        .get_int16_param       = int16ParamGetter,
        .get_int32_param       = int32ParamGetter,
        .get_int64_param       = int64ParamGetter,
        .get_uint8_param        = uint8ParamGetter,
        .get_uint16_param       = uint16ParamGetter,
        .get_uint32_param       = uint32ParamGetter,
        .get_uint64_param       = uint64ParamGetter,
        .get_double_param      = doubleParamGetter,
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
	"encoding/hex"
	"fmt"
	"runtime/cgo"
	"unsafe"

	"github.com/Gordy96/cgo_dl/dl"
	"github.com/Gordy96/evt-sim/modules/device"
)

//export goPacketDump
func goPacketDump(ctx *C.void, dir *C.cchar_t, data *C.void, size C.size_t) {
	c := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	if c != nil && c.log != nil {
		dat := unsafe.Slice((*byte)(unsafe.Pointer(data)), int(size))
		c.log(-1, C.GoString(dir)+":\n"+hex.Dump(dat))
	}
}

//export goLog
func goLog(ctx *C.void, level C.int, line *C.char) {
	c := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	if c != nil && c.log != nil {
		c.log(int(level), C.GoString(line))
	}
}

type portInterruptConfig struct {
	port string
	cb   C.interrupt_callback_t
}

//export attachPortInterrupt
func attachPortInterrupt(ctx *C.void, port *C.cchar_t, cb C.interrupt_callback_t) {
	c := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	p := C.GoString(port)
	cfg := portInterruptConfig{
		port: p,
		cb:   cb,
	}
	c.portInterrupts[p] = cfg
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

//export stringParamGetter
func stringParamGetter(ctx *C.void, name *C.cchar_t, buf *C.char, size C.size_t) C.int {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	istr, ok := a.params[C.GoString(name)]
	if !ok {
		return -1
	}

	str, ok := istr.(string)
	if !ok {
		return -1
	}

	src := []byte(str)

	ml := int(size) - 1
	if ml < 0 {
		return 0
	}

	n := len(src)
	if n > ml {
		n = ml
	}

	// C buffer as Go slice
	dst := unsafe.Slice((*byte)(unsafe.Pointer(buf)), int(size))

	// Copy bytes
	copy(dst[:n], src[:n])

	// Explicit null termination
	dst[n] = 0

	return C.int(n) // length excluding '\0'
}

//export int8ParamGetter
func int8ParamGetter(ctx *C.void, name *C.cchar_t, dst *C.int8_t) C.int {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	iint, ok := a.params[C.GoString(name)]
	if !ok {
		return -1
	}

	if i, ok := castAnyInt(iint); ok {
		*dst = C.int8_t(i)
		return 0
	}

	return -1
}

//export int16ParamGetter
func int16ParamGetter(ctx *C.void, name *C.cchar_t, dst *C.int16_t) C.int {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	iint, ok := a.params[C.GoString(name)]
	if !ok {
		return -1
	}

	if i, ok := castAnyInt(iint); ok {
		*dst = C.int16_t(i)
		return 0
	}

	return -1
}

//export int32ParamGetter
func int32ParamGetter(ctx *C.void, name *C.cchar_t, dst *C.int32_t) C.int {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	iint, ok := a.params[C.GoString(name)]
	if !ok {
		return -1
	}

	if i, ok := castAnyInt(iint); ok {
		*dst = C.int32_t(i)
		return 0
	}

	return -1
}

//export int64ParamGetter
func int64ParamGetter(ctx *C.void, name *C.cchar_t, dst *C.int64_t) C.int {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	iint, ok := a.params[C.GoString(name)]
	if !ok {
		return -1
	}

	if i, ok := castAnyInt(iint); ok {
		*dst = C.int64_t(i)
		return 0
	}

	return -1
}

//export uint8ParamGetter
func uint8ParamGetter(ctx *C.void, name *C.cchar_t, dst *C.uint8_t) C.int {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	iint, ok := a.params[C.GoString(name)]
	if !ok {
		return -1
	}

	if i, ok := castAnyInt(iint); ok {
		*dst = C.uint8_t(i)
		return 0
	}

	return -1
}

//export uint16ParamGetter
func uint16ParamGetter(ctx *C.void, name *C.cchar_t, dst *C.uint16_t) C.int {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	iint, ok := a.params[C.GoString(name)]
	if !ok {
		return -1
	}

	if i, ok := castAnyInt(iint); ok {
		*dst = C.uint16_t(i)
		return 0
	}

	return -1
}

//export uint32ParamGetter
func uint32ParamGetter(ctx *C.void, name *C.cchar_t, dst *C.uint32_t) C.int {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	iint, ok := a.params[C.GoString(name)]
	if !ok {
		return -1
	}

	if i, ok := castAnyInt(iint); ok {
		*dst = C.uint32_t(i)
		return 0
	}

	return -1
}

//export uint64ParamGetter
func uint64ParamGetter(ctx *C.void, name *C.cchar_t, dst *C.uint64_t) C.int {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	iint, ok := a.params[C.GoString(name)]
	if !ok {
		return -1
	}

	if i, ok := castAnyInt(iint); ok {
		*dst = C.uint64_t(i)
		return 0
	}

	return -1
}

func castAnyInt(i interface{}) (int64, bool) {
	if i, ok := i.(int64); ok {
		return i, true
	}
	if i, ok := i.(int32); ok {
		return int64(i), true
	}
	if i, ok := i.(int16); ok {
		return int64(i), true
	}
	if i, ok := i.(int8); ok {
		return int64(i), true
	}
	if i, ok := i.(int); ok {
		return int64(i), true
	}

	return 0, false
}

func castAnyUint(i interface{}) (uint64, bool) {
	if i, ok := i.(uint64); ok {
		return i, true
	}
	if i, ok := i.(uint32); ok {
		return uint64(i), true
	}
	if i, ok := i.(uint16); ok {
		return uint64(i), true
	}
	if i, ok := i.(uint8); ok {
		return uint64(i), true
	}
	if i, ok := i.(uint); ok {
		return uint64(i), true
	}

	return 0, false
}

//export doubleParamGetter
func doubleParamGetter(ctx *C.void, name *C.cchar_t, dst *C.double) C.int {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	idouble, ok := a.params[C.GoString(name)]
	if !ok {
		return -1
	}

	if d, ok := idouble.(float64); ok {
		*dst = C.double(d)
		return 0
	}
	if d, ok := idouble.(float32); ok {
		*dst = C.double(d)
		return 0
	}

	return -1
}

//export goRead
func goRead(ctx *C.void, port *C.cchar_t, buf *C.char, size C.size_t) C.int {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)

	if p, ok := a.ports[C.GoString(port)]; ok {
		temp := make([]byte, int(size))
		n, _ := p.Read(temp)

		C.memcpy(unsafe.Pointer(buf), unsafe.Pointer(&temp[0]), C.size_t(n))

		return C.int(n)
	}

	return 0
}

//export goWrite
func goWrite(ctx *C.void, port *C.cchar_t, buf *C.char, size C.size_t) C.int {
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
func dataGetter(ctx *C.void, name *C.cchar_t) *C.void {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	v, _ := a.mem[C.GoString(name)].(unsafe.Pointer)
	return (*C.void)(v)
}

//export dataSetter
func dataSetter(ctx *C.void, name *C.cchar_t, value *C.void) {
	a := cgo.Handle(unsafe.Pointer(ctx)).Value().(*Application)
	a.mem[C.GoString(name)] = unsafe.Pointer(value)
}

var _ device.Application = (*Application)(nil)

type Application struct {
	selfUnsafe      cgo.Handle
	initFunc        C.init_func_t
	shutdownFunc    C.shutdown_t
	portInterrupts  map[string]portInterruptConfig
	timerInterrupts map[string]timerInterruptConfig
	ports           map[string]device.Port
	mem             map[string]interface{}
	params          map[string]interface{}
	schedule        func(string, int)
	log             func(int, string)
	concurrent      bool
}

func (a *Application) Close() error {
	C.tShutdown(unsafe.Pointer(a.selfUnsafe), a.shutdownFunc)
	a.selfUnsafe.Delete()
	return nil
}

func (a *Application) Init(schedule func(string, int), ports ...device.Port) error {
	for _, port := range ports {
		a.ports[port.Name()] = port
	}
	a.schedule = schedule

	if a.concurrent {
		go C.tInit(unsafe.Pointer(a.selfUnsafe), a.initFunc)
	} else {
		C.tInit(unsafe.Pointer(a.selfUnsafe), a.initFunc)
	}

	return nil
}

func (a *Application) TriggerTimeInterrupt(key string) error {
	if i, ok := a.timerInterrupts[key]; ok {
		if a.concurrent {
			go C.tInterrupt(unsafe.Pointer(a.selfUnsafe), i.cb)
		} else {
			C.tInterrupt(unsafe.Pointer(a.selfUnsafe), i.cb)
		}
		if i.periodic {
			a.schedule(key, i.timeMS)
		}
	}

	return nil
}

func (a *Application) TriggerPortInterrupt(port string) error {
	if i, ok := a.portInterrupts[port]; ok {
		if a.concurrent {
			go C.tInterrupt(unsafe.Pointer(a.selfUnsafe), i.cb)
		} else {
			C.tInterrupt(unsafe.Pointer(a.selfUnsafe), i.cb)
		}
	}

	return nil
}

type Option func(*Application)

func WithLogger(logger func(int, string)) Option {
	return func(a *Application) {
		a.log = logger
	}
}

func WithParams(params map[string]interface{}) Option {
	return func(a *Application) {
		if a.params == nil {
			a.params = make(map[string]interface{})
		}
		a.params = params
	}
}

func WithParam[T any](name string, value T) Option {
	return func(a *Application) {
		if a.params == nil {
			a.params = make(map[string]interface{})
		}
		a.params[name] = value
	}
}

func WithConcurrency(c bool) Option {
	return func(a *Application) {
		a.concurrent = c
	}
}

func New(lib *dl.SO, opts ...Option) (*Application, error) {
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
		initFunc:        initFunc,
		shutdownFunc:    shutdownFunc,
		timerInterrupts: make(map[string]timerInterruptConfig),
		portInterrupts:  make(map[string]portInterruptConfig),
		ports:           make(map[string]device.Port),
		mem:             make(map[string]interface{}),
		params:          make(map[string]interface{}),
		concurrent:      false,
	}

	for _, opt := range opts {
		opt(a)
	}

	a.selfUnsafe = cgo.NewHandle(a)

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
