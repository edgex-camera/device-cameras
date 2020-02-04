module gitlab.jiangxingai.com/applications/edgex/device-service/device-cameras

go 1.13

replace github.com/edgexfoundry/device-sdk-go => ./internal/mod/device-sdk-go

replace github.com/edgexfoundry/go-mod-core-contracts => ./internal/mod/go-mod-core-contracts

replace github.com/edgexfoundry/go-mod-registry => ./internal/mod/go-mod-registry

require (
	github.com/edgexfoundry/device-sdk-go v1.0.0
	github.com/edgexfoundry/go-mod-core-contracts v0.1.29
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gorilla/mux v1.6.2
)
