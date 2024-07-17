# Golang iRacing SDK

Golang implementation of iRacing SDK.

This is a fork of the original [iracing-sdk](https://github.com/quimcalpe/iracing-sdk) package by [quimcalpe](https://github.com/quimcalpe) with some additional features and bug fixes.

## Install

1. Execute `go get github.com/riccardotornesello/irsdk-go`

## Usage

Simplest example:

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
    speed, _ := sdk.GetVar("Speed")
    fmt.Printf("Speed: %s", speed)
}
```

Get data in a loop live

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
        sdk.WaitForData(100 * time.Millisecond)
        speed, err := sdk.GetVar("Speed")
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Speed: %s", speed)
    }
}
```

Work with an offline ibt file

```go
reader, err := os.Open("data.ibt")
if err != nil {
    log.Fatal(err)
}
sdk = irsdk.Init(reader)
...
```

## Examples

- [Export](examples/export) Telemetry Data and Session yaml to files

- Broadcast [Commands](examples/commands) to iRacing

- Simple [Dashboard](examples/dashboard) for external monitors or phones
