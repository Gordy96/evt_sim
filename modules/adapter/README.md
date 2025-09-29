# Adapter
Adapter package is meant to implement "software" module for simulation Component by loading shared objects/dll's 
and exposing generic interface for time and pin based interrupts (for MCU like device components)
as well as "reading" and "writing" to ports/gates.

Use `include/plugin.h` for type definitions. Adapter first initializes library itself by calling `init_lib` 
function with go callbacks packed in `interface_t` struct, and then `init` and `shutdown` for each instance, 
so library is instantiated once per simulation (cached) but components can spawn 
multiple "instances" of `Adapter`

## warning
Developers must pay close attention to what C code allocates/deallocates. `Adapter` does not control this memory 
aside from calling `shutdown`

## warning 2
`void *ctx` is a "cgo pointer" (global sync.Map of pointers) and should not be meddled with at any point in SO/DLL (C)

## Examples
Sample plugin
```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "plugin.h"

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
    int n = ENV.read_port(ctx, "port", buf);
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
    ENV.attach_pin_interrupt(ctx, 2, cb);
}

void shutdown(void *ctx) {
    ENV.log("shutting down");
    App *app = ENV.data_getter(ctx, "this");
    free(app);
}
```
