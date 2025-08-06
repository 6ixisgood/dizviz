// +build noled

package main

// Stub config for testing without LED hardware
type DefaultConfig struct {
	Rows                    int
	Cols                    int
	Parallel                int
	ChainLength             int
	Brightness              int
	HardwareMapping         string
	ShowRefreshRate         bool
	InverseColors           bool
	DisableHardwarePulsing  bool
	GpioSlowdown            int
	RateLimitHz             int
}
