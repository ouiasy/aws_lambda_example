.PHONY: compile
compile:
	GOOS=linux GOARCH=amd64 go build -o ./output/bootstrap main.go
	zip -j ./output/main.zip ./output/bootstrap