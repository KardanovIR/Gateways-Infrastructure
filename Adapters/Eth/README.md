# GatewaysInfrastructure

Run :

    go run main.go [parameters...] 
    
    Parameters:
       -config-path     path to config, by default "./config/config.yml"

Recompile protobuf :
    
    protoc -I grpc/ grpc/ethAdapter.proto --go_out=plugins=grpc:grpc
