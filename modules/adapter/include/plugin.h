#ifndef PLUGIN_H
#define PLUGIN_H

typedef int   (*read_cb_t)(void *ctx, char* port, char* buf);
typedef int   (*write_cb_t)(void *ctx, char* port, char* buf, int size);
typedef void  (*interrupt_callback_t)(void *ctx);
typedef void  (*attach_pin_interrupt_t)(void *ctx, int pin, interrupt_callback_t cb);
typedef void  (*attach_time_interrupt_t)(void *ctx, int time_ms, short periodic, interrupt_callback_t cb);
typedef void  (*shutdown_t)(void *ctx);
typedef void* (*getter_t)(void *ctx, char* name);
typedef void  (*setter_t)(void *ctx, char* name, void* value);
typedef void  (*log_t)(char *line);

typedef struct {
    read_cb_t               read_port;
    write_cb_t              write_port;
    attach_pin_interrupt_t  attach_pin_interrupt;
    attach_time_interrupt_t attach_time_interrupt;
    getter_t                data_getter;
    setter_t                data_setter;
    log_t                   log;
} interface_t;

typedef void (*lib_init_func_t)(interface_t iface);
typedef void (*init_func_t)(void *ctx);

#endif