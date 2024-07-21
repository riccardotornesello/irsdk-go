# Golang iRacing SDK

Golang implementation of iRacing SDK.

This project is a functional prototype and is not yet ready for production use.

It can be used to read telemetry and session data but some features are still missing.

## Install

1. Execute `go get github.com/riccardotornesello/irsdk-go`
2. Enjoy

## Usage

Get session info:

```go
package main

import (
    "fmt"
    "github.com/riccardotornesello/irsdk-go"
)

func main() {
    var sdk irsdk.IRSDK
    sdk = irsdk.Init(nil)
    defer sdk.Close()

    userId := sdk.Session.DriverInfo.DriverUserID
    fmt.Printf("User ID: %s", userId)
}
```

Get telemetry in a loop live

```go
package main

import (
    "fmt"
    "log"

    "github.com/riccardotornesello/irsdk-go"
)

func main() {
    var sdk irsdk.IRSDK
    sdk = irsdk.Init(nil)
    defer sdk.Close()

    for {
        sdk.Update(true)
        speed := sdk.Telemetry["Speed"]
        fmt.Printf("Speed: %s", speed)
    }
}
```

## Examples

- [Export](examples/export) Telemetry Data and Session yaml to files

- Broadcast [Commands](examples/commands) to iRacing

- Simple [Dashboard](examples/dashboard) for external monitors or phones

- [Timing Console](examples/times) to show live timing data

## Missing features

- [ ] Sending commands to iRacing
- [ ] Better documentation

## Credits

- [quimcalpe's iracing-sdk](https://github.com/quimcalpe/iracing-sdk) package for the original package, used to better understand the iRacing SDK
- Iracing C++ SDK for the original implementation of the SDK, useful to get all the defines and structs
