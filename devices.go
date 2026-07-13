package main

import "github.com/albionmarketgg/dataclient/internal/capture"

func captureDevices() ([]string, error) { return capture.Devices() }
