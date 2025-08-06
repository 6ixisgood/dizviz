// +build !noled

package main

import "github.com/sixisgoood/go-rpi-rgb-led-matrix"

// Use the real matrix config when not using stub
type DefaultConfig = rgbmatrix.HardwareConfig
