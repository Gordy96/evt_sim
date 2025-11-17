#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "../include/plugin.h"

interface_t ENV = {0};

void init_lib(interface_t iface) {
    ENV = iface;
}

typedef struct {
    int counter;
} App;

void cb(void *ctx) {
    App *app = ENV.data_getter(ctx, "this");

    char buf[256];
    int n = ENV.read_port(ctx, "port", buf, 256);
    if (n > 0) {
        app->counter += 1;
        buf[n] = '\0';
        char temp[256];
        n = sprintf(temp, "%s %d", buf, app->counter);
        ENV.write_port(ctx, "port", temp, n); // echo back
    }
}

void init(void *ctx) {
    App *app = malloc(sizeof(App));
    app->counter = 0;
    ENV.data_setter(ctx, "this", app);
    ENV.attach_port_interrupt(ctx, "port", cb);
}

void shutdown(void *ctx) {
    ENV.log("shutting down");
    App *app = ENV.data_getter(ctx, "this");
    free(app);
}