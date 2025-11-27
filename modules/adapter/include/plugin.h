#ifndef PLUGIN_H
#define PLUGIN_H

typedef int    (*read_cb_t)(void *ctx, char* port, char* buf, int size);
typedef int    (*write_cb_t)(void *ctx, char* port, char* buf, int size);
typedef void   (*interrupt_callback_t)(void *ctx);
typedef void   (*attach_port_interrupt_t)(void *ctx, char* port, interrupt_callback_t cb);
typedef void   (*attach_time_interrupt_t)(void *ctx, int time_ms, short periodic, interrupt_callback_t cb);
typedef void   (*shutdown_t)(void *ctx);
typedef void*  (*getter_t)(void *ctx, char* name);
typedef void   (*setter_t)(void *ctx, char* name, void* value);
typedef int    (*string_param_getter_t)(void *ctx, char* name, char* buf, int size);
typedef int    (*int_param_getter_t)(void *ctx, char* name, int* dst);
typedef int    (*double_param_getter_t)(void *ctx, char* name, double* dst);
typedef void   (*log_t)(void *ctx, char *line);

typedef struct {
    read_cb_t               read_port;
    write_cb_t              write_port;
    attach_port_interrupt_t attach_port_interrupt;
    attach_time_interrupt_t attach_time_interrupt;
    getter_t                get_data;
    setter_t                set_data;
    string_param_getter_t   get_string_param;
    int_param_getter_t      get_int_param;
    double_param_getter_t   get_double_param;
    log_t                   log;
} interface_t;

typedef void (*lib_init_func_t)(interface_t iface);
typedef void (*init_func_t)(void *ctx);

#endif