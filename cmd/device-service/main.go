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

	"github.com/edgexfoundry/device-sdk-go"
	"github.com/edgexfoundry/device-sdk-go/pkg/jxstartup"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/driver"
	"gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras/internal/restful"
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
	err := jxstartup.StartServiceWithHandler(serviceName, device.Version, &driver.CurrentDriver, restful.InitRestRoutes, staticPath)
	if err != nil {
		panic(err)
	}
}
