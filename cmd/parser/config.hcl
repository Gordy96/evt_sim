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

variable "devices" {
  type = list(object({
    id = string
    position = object({
      lat  = float64
      lon  = float64
      elev = float64
    })
    radio = object({
      frequency_hz = int64
      power        = int64
      fade_margin  = int64
    })
    app = object({
      so_path = string
      params  = object({
        address   = uint32
        dest      = uint32
        parallel  = bool
        initiator = bool
      })
    })
  }))
  default = [
    {
      id = "first"
      radio = default_radio
      position = {
        lat  = 50.45624
        lon  = 30.36545
        elev = 200
      }
      app = {
        so_path = "${APP_SO_PATH}"
        params = {
          address   = 1
          dest      = 3
          parallel  = true
          initiator = true
        }
      }
    },
    {
      id = "second"
      radio = default_radio
      position = {
        lat  = 50.45422
        lon  = 30.44862
        elev = 200
      }
      app = {
        so_path = "${APP_SO_PATH}"
        params = {
          address   = 2
          parallel  = true
        }
      }
    },
    {
      id = "third"
      radio = default_radio
      position = {
        lat  = 50.44812
        lon  = 30.525
        elev = 200
      }
      app = {
        so_path = "${APP_SO_PATH}"
        params = {
          address   = 3
          parallel  = true
        }
      }
    }
  ]
}

# radio_medium {
#   # background_noise_level = 10
# }

module embedded {
  for_each = devices
  id = each.id
  position WGS84 {
    lat  = each.position.lat
    lon  = each.position.lon
    elev = each.position.elev
  }
  radio lora {
    id = "port"
    frequency_hz = each.radio.frequency_hz
    power     = each.radio.power
    fade_margin = each.radio.fade_margin
  }
  application shared {
    path      = each.app.so_path
    dump_packets = true
    parameters = each.app.params
  }
}
