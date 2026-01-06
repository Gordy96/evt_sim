name = "test"

# radio_medium {
#   # background_noise_level = 10
# }

module embedded first {
  position WGS84 {
    lat  = 50.45624
    lon  = 30.36545
    elev = 200
  }
  radio lora {
    id = "port"
    frequency_hz = 433000000
    power     = 20
    fade_margin = 10000 #10Km
  }
  application shared {
    path      = "${APP_SO_PATH}"
    parallel  = 1
    address   = 1
    initiator = 1
    dest      = 3
    dump_packets = true
  }
}

module embedded second {
  position WGS84 {
    lat  = 50.45422
    lon  = 30.44862
    elev = 200
  }
  radio lora {
    id = "port"
    frequency_hz = 433000000
    power     = 20
    fade_margin = 10000 #10Km
  }
  application shared {
    path      = "${APP_SO_PATH}"
    parallel  = 1
    address   = 2
    dump_packets = true
  }
}

module embedded third {
  position WGS84 {
    lat  = 50.44812
    lon  = 30.525
    elev = 200
  }
  radio lora {
    id = "port"
    frequency_hz = 433000000
    power     = 20
    fade_margin = 10000 #10Km
  }
  application shared {
    path      = "${APP_SO_PATH}"
    parallel  = 1
    address   = 3
    dump_packets = true
  }
}
