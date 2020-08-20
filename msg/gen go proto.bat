protoc.exe --plugin=protoc-gen-go=protoc-gen-go.exe --go_out=. msg.proto

del msg.go
ren msg.pb.go msg.go

pause;