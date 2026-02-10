name = "test"

variable "default_radio" {
  type = object({
    frequency_hz = int64
    power        = int64
    fade_margin  = int64
  })
  default = {
    frequency_hz = 433000000
    power        = 20
    fade_margin  = 10000 #10Km
  }
}

variable "routes" {
  type = list(string)
  default = ["3:2|2:2|4:4", "1:1|3:3|4:4", "1:2|2:2|4:4", "1:1|3:3|2:2"]
}

variable "devices" {
  type = list(object({
    lat  = float64
    lon  = float64
    elev = float64
  }))
  default = [
    {
      lat  = 50.45624
      lon  = 30.36545
      elev = 200
    },
    {
      lat  = 50.45422
      lon  = 30.44862
      elev = 200
    },
    {
      lat  = 50.44812
      lon  = 30.525
      elev = 200
    },
    {
      lat  = 50.44505
      lon  = 30.44371
      elev = 200
    }
  ]
}

# radio_medium {
#   # background_noise_level = 10
# }
realtime = false

module embedded {
  for_each = devices
  id = itoa(iter + 1)
  position WGS84 {
    lat  = each.lat
    lon  = each.lon
    elev = each.elev
  }
  radio lora {
    id = "port"
    frequency_hz = default_radio.frequency_hz
    power        = default_radio.power
    fade_margin  = default_radio.fade_margin
  }
  application shared {
    path         = "${APP_SO_PATH}"
    dump_packets = true
    parameters   = {
      address   = iter + 1
      dest      = iter == 0 ? 3 : 0
      parallel  = true
      initiator = iter == 0
      routing_table = routes[iter]
    }
  }
}
