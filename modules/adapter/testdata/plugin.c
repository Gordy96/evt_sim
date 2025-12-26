#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include "../include/plugin.h"

interface_t ENV = {0};

void init_lib(interface_t iface) {
    ENV = iface;
}

typedef struct {
    uint8_t is_initiator;
    int32_t counter;
    double  factor;
    char    name[256];
} App;

void cb(void *ctx) {
    App *app = ENV.get_data(ctx, "this");

    char buf[256];
    int n = ENV.read_port(ctx, "port", buf, 256);
    if (n > 0) {
        ENV.dump_packet(ctx, "IN", buf, n);
        buf[n] = '\0';
        ENV.log(ctx, DebugLevel, "received message");
        ENV.log(ctx, DebugLevel, buf);
        if (app->counter > 0) {
            app->counter -= 1;
            char temp[256];
            n = sprintf(temp, "%s %f %s %d", app->name, app->factor, buf, app->counter);
            ENV.write_port(ctx, "port", temp, n); // echo back
            ENV.dump_packet(ctx, "OUT", temp, n);
        }
    }
}

void init(void *ctx) {
    App *app = malloc(sizeof(App));
    ENV.get_int32_param(ctx, "counter", &app->counter);
    ENV.get_uint8_param(ctx, "initiator", &app->is_initiator);
    ENV.get_double_param(ctx, "factor", &app->factor);
    int n = ENV.get_string_param(ctx, "name", app->name, 256);
    ENV.set_data(ctx, "this", app);
    ENV.attach_port_interrupt(ctx, "port", cb);
    ENV.log(ctx, InfoLevel, "initializing");

    if (app->is_initiator) {
        char temp[256];
        int n = sprintf(temp, "%s says hello", app->name);
        ENV.log(ctx, DebugLevel, "sending initial message");
        ENV.log(ctx, DebugLevel, temp);
        ENV.write_port(ctx, "port", temp, n);
        ENV.dump_packet(ctx, "OUT", temp, n);
    }
}

void shutdown(void *ctx) {
    ENV.log(ctx, InfoLevel, "shutting down");
    App *app = ENV.get_data(ctx, "this");
    free(app);
}