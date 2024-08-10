protoc --proto_path=. --proto_path=./proto --go_out=. packet.proto
protoc --proto_path=. --proto_path=./proto --csharp_out=./cli packet.proto