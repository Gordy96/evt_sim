name = "test"

# radio_medium {
#   # background_noise_level = 10
# }

module embedded first {
  position WGS84 {
    lon  = 50.451617
    lat  = 30.300575
    elev = 200
  }
  radio lora {
    id = "port"
    frequency_hz = 433000000
    power     = 20
    fade_margin = 10000 #10Km
  }
  application shared {
    path      = "modules/adapter/testdata/plugin.so"
    factor    = 1.2
    counter   = 0
    name      = "first"
    initiator = 1
  }
}

module embedded second {
  position WGS84 {
    lon  = 50.462015
    lat  = 30.355995
    elev = 200
  }
  radio lora {
    id = "port"
    frequency_hz = 433000000
    power     = 20
    fade_margin = 10000 #10Km
  }
  application shared {
    path    = "modules/adapter/testdata/plugin.so"
    factor  = 3.4
    counter = 1
    name    = "second"
  }
}
