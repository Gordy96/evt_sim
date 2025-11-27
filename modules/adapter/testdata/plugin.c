#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "../include/plugin.h"

interface_t ENV = {0};

void init_lib(interface_t iface) {
    ENV = iface;
}

typedef struct {
    int    counter;
    double factor;
    char   name[256];
} App;

void cb(void *ctx) {
    App *app = ENV.get_data(ctx, "this");

    char buf[256];
    int n = ENV.read_port(ctx, "port", buf, 256);
    if (n > 0) {
        app->counter += 1;
        buf[n] = '\0';
        char temp[256];
        n = sprintf(temp, "%s %f %s %d", app->name, app->factor, buf, app->counter);
        ENV.write_port(ctx, "port", temp, n); // echo back
    }
}

void init(void *ctx) {
    App *app = malloc(sizeof(App));
    ENV.get_int_param(ctx, "counter", &app->counter);
    ENV.get_double_param(ctx, "factor", &app->factor);
    ENV.get_string_param(ctx, "name", app->name, 256);
    ENV.set_data(ctx, "this", app);
    ENV.attach_port_interrupt(ctx, "port", cb);
}

void shutdown(void *ctx) {
    ENV.log(ctx, "shutting down");
    App *app = ENV.get_data(ctx, "this");
    free(app);
}