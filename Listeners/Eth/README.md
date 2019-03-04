# GatewaysInfrastructure

Run :

    go run main.go [parameters...] 
    
    Parameters:
       -config-path     path to config, by default "./config/config.yml"
       -config-path=./Listeners/Eth/config/config.yml for debug in JetBrains Goland
       
Recompile protobuf :
    
    protoc -I grpc/ grpc/eth_listener.proto --go_out=plugins=grpc:grpc

Docker :
    
    Build container:
    
    docker build --rm -t gateways-eth-listener:latest .
    
    Run container
    docker run --rm -d -p 5001:5001 --env-file=config/dev.env --name gateways-eth-listener gateways-eth-listener:latest
    
    Read logs from container
    docker logs gateways-eth-listener