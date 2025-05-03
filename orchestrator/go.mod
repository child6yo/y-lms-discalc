module github.com/child6yo/y-lms-discalc/orchestrator

go 1.23.0

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/child6yo/y-lms-discalc/shared v0.0.0
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/mattn/go-sqlite3 v1.14.28
	google.golang.org/grpc v1.72.0
)

require (
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/child6yo/y-lms-discalc/shared => ../shared
