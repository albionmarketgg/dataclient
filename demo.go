package main

import "github.com/niick1231/albionmarket_dataclient/internal/capture"

func captureDevices() ([]string, error) { return capture.Devices() }
