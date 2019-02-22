# GatewaysInfrastructure

Run :

    go run main.go [parameters...] 
    
    Parameters:
       -config-path     path to config, by default "./config/config.yml"
       -debug           set debug mode. It set log level to Debug and log format (plain text instead of json)

Recompile protobuf :
    
    protoc -I grpc/ grpc/ethAdapter.proto --go_out=plugins=grpc:grpc
