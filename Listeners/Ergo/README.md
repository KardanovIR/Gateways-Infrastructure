# GatewaysInfrastructure

Run :

    go run main.go [parameters...] 
    
    Parameters:
       -config-path     path to config, by default "./config/config.yml"
       -config-path=./Listeners/Ergo/config/config.yml for debug in JetBrains Goland
       
Recompile protobuf :
    
    protoc -I grpc/ grpc/ergo_listener.proto --go_out=plugins=grpc:grpc
    protoc -I grpc/client grpc/client/callback_service.proto --go_out=plugins=grpc:grpc/client

Docker :
    
    Build container:
    
    docker build --rm -t gateways-ergo-listener:latest .
    
    Run container
    docker run --rm -d -p 5001:5001 --env-file=config/dev.env --name gateways-ergo-listener gateways-ergo-listener:latest
    
    Read logs from container
    docker logs gateways-ergo-listener