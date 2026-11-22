$ErrorActionPreference = 'Stop'

# Выполнить из корневого каталога проекта в PowerShell.
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.2
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

# protoc и каталог Go/bin должны быть доступны через PATH.
protoc --proto_path=. --go_out=. --go_opt=module=example.com/student-service --go-grpc_out=. --go-grpc_opt=module=example.com/student-service proto/student.proto
