## Usage

> go run ./server/server.go [-host=] [-port=] &
> go run ./client/client.go [-host=] [-port=] [-backhost=] [-backport=]


## Create a new command

1. Add the method to the CommandRunner service in proto/server.proto
2. Generate the go code by running codegen.sh
3. Implement the method in server.go
4. Add the call to the switch in client.go
