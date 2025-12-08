#ifndef PLUGIN_H
#define PLUGIN_H

#include <stdint.h>

typedef enum {
  	DebugLevel = - 1,
  	// InfoLevel is the default logging priority.
  	InfoLevel,
  	// WarnLevel logs are more important than Info, but don't need individual
  	// human review.
  	WarnLevel,
  	// ErrorLevel logs are high-priority. If an application is running smoothly,
  	// it shouldn't generate any error-level logs.
  	ErrorLevel,
  	// DPanicLevel logs are particularly important errors. In development the
  	// logger panics after writing the message.
  	DPanicLevel,
  	// PanicLevel logs a message, then panics.
  	PanicLevel,
  	// FatalLevel logs a message, then calls os.Exit(1).
  	FatalLevel
} LogLevel;

typedef int    (*read_cb_t)(void *ctx, char* port, char* buf, int size);
typedef int    (*write_cb_t)(void *ctx, char* port, char* buf, int size);
typedef void   (*interrupt_callback_t)(void *ctx);
typedef void   (*attach_port_interrupt_t)(void *ctx, char* port, interrupt_callback_t cb);
typedef void   (*attach_time_interrupt_t)(void *ctx, int time_ms, short periodic, interrupt_callback_t cb);
typedef void   (*shutdown_t)(void *ctx);
typedef void*  (*getter_t)(void *ctx, char* name);
typedef void   (*setter_t)(void *ctx, char* name, void* value);
typedef int    (*string_param_getter_t)(void *ctx, char* name, char* buf, int size);
typedef int    (*int8_param_getter_t)(void *ctx, char* name, int8_t* dst);
typedef int    (*int16_param_getter_t)(void *ctx, char* name, int16_t* dst);
typedef int    (*int32_param_getter_t)(void *ctx, char* name, int32_t* dst);
typedef int    (*int64_param_getter_t)(void *ctx, char* name, int64_t* dst);
typedef int    (*uint8_param_getter_t)(void *ctx, char* name, uint8_t* dst);
typedef int    (*uint16_param_getter_t)(void *ctx, char* name, uint16_t* dst);
typedef int    (*uint32_param_getter_t)(void *ctx, char* name, uint32_t* dst);
typedef int    (*uint64_param_getter_t)(void *ctx, char* name, uint64_t* dst);
typedef int    (*double_param_getter_t)(void *ctx, char* name, double* dst);
typedef void   (*log_t)(void *ctx, LogLevel level, char *line);

typedef struct {
    read_cb_t                read_port;
    write_cb_t               write_port;
    attach_port_interrupt_t  attach_port_interrupt;
    attach_time_interrupt_t  attach_time_interrupt;
    getter_t                 get_data;
    setter_t                 set_data;
    string_param_getter_t    get_string_param;
    int8_param_getter_t      get_int8_param;
    int16_param_getter_t     get_int16_param;
    int32_param_getter_t     get_int32_param;
	int64_param_getter_t     get_int64_param;
    uint8_param_getter_t     get_uint8_param;
    uint16_param_getter_t    get_uint16_param;
    uint32_param_getter_t    get_uint32_param;
	uint64_param_getter_t    get_uint64_param;
    double_param_getter_t    get_double_param;
    log_t                    log;
} interface_t;

typedef void (*lib_init_func_t)(interface_t iface);
typedef void (*init_func_t)(void *ctx);

#endif