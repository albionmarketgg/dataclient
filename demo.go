package main

import "github.com/albionmarketgg/data-client/internal/capture"

func captureDevices() ([]string, error) { return capture.Devices() }
