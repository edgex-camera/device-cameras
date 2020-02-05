module gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras

go 1.13

replace github.com/edgexfoundry/device-sdk-go => ./internal/mod/device-sdk-go

replace github.com/edgexfoundry/go-mod-core-contracts => ./internal/mod/go-mod-core-contracts

replace github.com/edgexfoundry/go-mod-registry => ./internal/mod/go-mod-registry

require (
	github.com/beevik/etree v1.1.0 // indirect
	github.com/edgexfoundry/device-sdk-go v1.0.0
	github.com/edgexfoundry/go-mod-core-contracts v0.1.29
	github.com/elgs/gostrgen v0.0.0-20161222160715-9d61ae07eeae // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gorilla/mux v1.6.2
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/yakovlevdmv/Golang-iso8601-duration v0.0.0-20180403125811-e5db0413b903 // indirect
	github.com/yakovlevdmv/WS-Discovery v0.0.0-20180512141937-16170c6c3677 // indirect
	github.com/yakovlevdmv/goonvif v0.0.0-20180517145634-8181eb3ef2fb
	github.com/yakovlevdmv/gosoap v0.0.0-20180512142237-299a954b1c6d // indirect
	gopkg.in/evanphx/json-patch.v4 v4.5.0
)

replace github.com/satori/go.uuid v1.2.0 => github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
