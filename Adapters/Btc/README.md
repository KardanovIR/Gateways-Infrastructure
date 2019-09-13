# GatewaysInfrastructure

Run :

    go run main.go [parameters...] 
    
    Parameters:
       -config-path     path to config, by default "./config/config.yml"
       -debug           set debug mode. It set log level to Debug and log format (plain text instead of json)
       
       -config-path=./Adapters/Btc/config/config.yml

Recompile protobuf :
    
    protoc -I grpc/ grpc/btc_adapter.proto --go_out=plugins=grpc:grpc
    
Docker :
    
    Build container:
    
    docker build --rm -t gateways-btc-adapter:latest .
    
    Run container
    docker run --rm -d -p 5001:5001 --env-file=config/dev.env --name gateways-btc-adapter gateways-btc-adapter:latest
    
    Read logs from container
    docker logs gateways-btc-adapter
