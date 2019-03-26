# GatewaysInfrastructure: Router

Run :

    go run main.go [parameters...] 
    
    Parameters:
       -config-path     path to config, by default "./config/config.yml"
       -config-path=./Router/config/config.yml for debug in JetBrains Goland
       
       -debug           set debug mode. It set log level to Debug and log format (plain text instead of json)
       
Recompile protobuf :
    
    router server api:
    protoc -I grpc/ grpc/router.proto --go_out=plugins=grpc:grpc
    
    protobuf client for eth adapter:
    protoc --proto_path=../Adapters/Eth/grpc --go_out=plugins=grpc:grpc/ethAdapter ./../Adapters/Eth/grpc/eth_adapter.proto
 
    protobuf client for waves adapter:
    protoc --proto_path=../Adapters/Waves/grpc --go_out=plugins=grpc:grpc/wavesAdapter ./../Adapters/Waves/grpc/waves_adapter.proto
    
    protobuf client for eth listener:
    protoc --proto_path=../Listeners/Core/grpc --go_out=plugins=grpc:grpc/ethListener ./../Listeners/Core/grpc/eth_listener.proto
    
    protobuf client for waves listener:
    protoc --proto_path=../Listeners/Waves/grpc --go_out=plugins=grpc:grpc/wavesListener ./../Listeners/Waves/grpc/waves_listener.proto
