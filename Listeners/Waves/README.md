# GatewaysInfrastructure

Run :

    go run main.go [parameters...] 
    
    Parameters:
       -config-path     path to config, by default "./config/config.yml"
       -config-path=./Listeners/Waves/config/config.yml for debug in JetBrains Goland
       
Recompile protobuf :
    
    protoc -I grpc/ grpc/waves_listener.proto --go_out=plugins=grpc:grpc
    protoc -I grpc/client grpc/client/callback_service.proto --go_out=plugins=grpc:grpc/client

Docker :
    
    Build container:
    
    docker build --rm -t gateways-waves-listener:latest .
    
    Run container
    docker run --rm -d -p 5001:5001 --env-file=config/dev.env --name gateways-waves-listener gateways-waves-listener:latest
    
    Read logs from container
    docker logs gateways-waves-listener