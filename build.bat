go run -tags=create_router app.go emptyLoadRouteHandlers.go
go test -count=1
go build
