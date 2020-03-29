module github.com/edgex-camera/device-cameras

go 1.13

replace github.com/edgex-camera/device-sdk-go => ./internal/mod/device-sdk-go

replace github.com/satori/go.uuid v1.2.0 => github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b

replace github.com/edgexfoundry/go-mod-core-contracts => github.com/edgexfoundry/go-mod-core-contracts v0.1.14

replace github.com/edgexfoundry/go-mod-registry => github.com/edgexfoundry/go-mod-registry v0.1.9

require (
	github.com/beevik/etree v1.1.0 // indirect
	github.com/edgex-camera/device-sdk-go v1.1.2
	github.com/edgexfoundry/go-mod-core-contracts v0.1.52
	github.com/elgs/gostrgen v0.0.0-20161222160715-9d61ae07eeae // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gorilla/mux v1.7.4
	github.com/satori/go.uuid v1.2.0
	github.com/yakovlevdmv/Golang-iso8601-duration v0.0.0-20180403125811-e5db0413b903 // indirect
	github.com/yakovlevdmv/WS-Discovery v0.0.0-20180512141937-16170c6c3677 // indirect
	github.com/yakovlevdmv/goonvif v0.0.0-20180517145634-8181eb3ef2fb
	github.com/yakovlevdmv/gosoap v0.0.0-20180512142237-299a954b1c6d // indirect
	gopkg.in/evanphx/json-patch.v4 v4.5.0
)
