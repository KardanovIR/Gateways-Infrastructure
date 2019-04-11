# GatewaysInfrastructure: Router

Run :

    go run main.go [parameters...] 
    
    Parameters:
       -config-path     path to config, by default "./config/config.yml"
       -config-path=./Router/config/config.yml for debug in JetBrains Goland
       
       -debug           set debug mode. It set log level to Debug and log format (plain text instead of json)
       
Recompile protobuf :
    
    protobuf client for all adapters and listeners
    protoc -I grpc/blockchain grpc/blockchain/blockchain_services.proto  --go_out=plugins=grpc:grpc/blockchain
   
    router server api:
    protoc -I grpc/ grpc/router.proto --go_out=plugins=grpc:grpc
    
Docker :
    
    Build container:
    
    docker build --rm -t gateways-router:latest .
    
    Run container
    docker run --rm -d -p 5001:5001 --env-file=config/local.env --name gateways-router gateways-router:latest
    
    Read logs from container
    docker logs gateways-router

Nginx :
    
    nginx.conf shoul be placed to /usr/local/etc/nginx/nginx.conf
    start: sudo nginx
    stop:  sudo nginx -s quit
    