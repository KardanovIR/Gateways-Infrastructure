# GatewaysInfrastructure

Run :

    go run main.go [parameters...] 
    
    Parameters:
       -config-path     path to config, by default "./config/config.yml"
       -debug           set debug mode. It set log level to Debug and log format (plain text instead of json)

Recompile protobuf :
    
    protoc -I grpc/ grpc/ergo_adapter.proto --go_out=plugins=grpc:grpc
    
Docker :
    
    Build container:
    
    docker build --rm -t gateways-ergo-adapter:latest .
    
    Run container
    docker run --rm -d -p 5001:5001 --env-file=config/local.env --name gateways-ergo-adapter gateways-ergo-adapter:latest
    
    Read logs from container
    docker logs gateways-ergo-adapter