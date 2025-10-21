.PHONY: compile-get_all_weather
compile-get_all_weather:
	GOOS=linux GOARCH=amd64 go build -o lambda_functions/get_all_weather/output/bootstrap lambda_functions/get_all_weather/main.go
	zip -j lambda_functions/get_all_weather/output/main.zip lambda_functions/get_all_weather/output/bootstrap

.PHONY: compile-random_weather
compile-random_weather:
	GOOS=linux GOARCH=amd64 go build -o lambda_functions/random_weather/output/bootstrap lambda_functions/random_weather/main.go
	zip -j lambda_functions/random_weather/output/main.zip lambda_functions/random_weather/output/bootstrap

.PHONY: compile-update_weather_from_csv
compile-update_weather_from_csv:
	GOOS=linux GOARCH=amd64 go build -o lambda_functions/update_weather_from_csv/output/bootstrap lambda_functions/update_weather_from_csv/main.go
	zip -j lambda_functions/update_weather_from_csv/output/main.zip lambda_functions/update_weather_from_csv/output/bootstrap