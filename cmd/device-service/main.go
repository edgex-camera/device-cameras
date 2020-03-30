// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a device service of cameras.
package main

import (
	"flag"

	"github.com/edgex-camera/device-cameras/internal/driver"
	"github.com/edgex-camera/device-cameras/internal/restful"
	"github.com/edgex-camera/device-sdk-go"
	"github.com/edgex-camera/device-sdk-go/pkg/camstartup"
)

const (
	serviceName = "device-cameras"
	staticPath  = "/app/frontend"
)

func main() {
	var processMethod string
	flag.StringVar(&processMethod, "pm", "ffmpeg", "which process should be used to process video, eg. ffmpeg, gst-launch-1.0")
	flag.Parse()

	driver.CurrentDriver = driver.Driver{ProcessMethod: processMethod}
	err := camstartup.StartServiceWithHandler(serviceName, device.Version, &driver.CurrentDriver, restful.InitRestRoutes, staticPath)
	if err != nil {
		panic(err)
	}
}
