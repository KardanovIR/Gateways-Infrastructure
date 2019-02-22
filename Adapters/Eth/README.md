# GatewaysInfrastructure

Run :

    go run main.go [parameters...] 
    
    Parameters:
       -config-path     path to config, by default "./config/config.yml"

Recompile protobuf :
    
    protoc -I grpc/ grpc/ethAdapter.proto --go_out=plugins=grpc:grpc
    
Docker :
    
    Build container:
    
    docker build --rm -t gateways-eth-adapter:latest .
    
    Run container
    docker run --rm -d -p 5001:5001 --env-file=config/dev.env --name gateways-eth-adapter gateways-eth-adapter:latest
    
    Read logs from container
    docker logs gateways-eth-adapter
