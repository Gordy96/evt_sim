name = "test"

module embedded first {
  radio lora {
    frequency = 433.0
    power     = 20
  }
  application shared {
    path = "modules/adapter/testdata/plugin.so"
    factor = 1.2
    counter = 10
    name = "first"
  }
}

module embedded second {
  radio lora {
    frequency = 433.0
    power     = 20
  }
  application shared {
    path = "modules/adapter/testdata/plugin.so"
    factor = 3.4
    counter = 20
    name = "second"
  }
}
