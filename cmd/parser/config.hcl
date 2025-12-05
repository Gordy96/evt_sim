name = "test"

# radio_medium {
#   # background_noise_level = 10
# }

module embedded first {
  position merkator {
    lon  = 100.00
    lat  = 222.00
    elev = 200
  }
  radio lora {
    frequency = 433.0
    power     = 20
    fade_margin = 10000 #10Km
  }
  application shared {
    path    = "modules/adapter/testdata/plugin.so"
    factor  = 1.2
    counter = 10
    name    = "first"
  }
}

module embedded second {
  position merkator {
    lon  = 100.00
    lat  = 222.00
    elev = 200
  }
  radio lora {
    frequency = 433.0
    power     = 20
    fade_margin = 10000 #10Km
  }
  application shared {
    path    = "modules/adapter/testdata/plugin.so"
    factor  = 3.4
    counter = 20
    name    = "second"
  }
}
