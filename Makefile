.PHONY: compile-get_all_weather
compile-get_all_weather:
	mkdir -p bin/get_all_weather
	GOOS=linux GOARCH=amd64 go build -o bin/get_all_weather/bootstrap lambda_functions/get_all_weather/main.go
	zip -j bin/get_all_weather/main.zip bin/get_all_weather/bootstrap

.PHONY: compile-random_weather
compile-random_weather:
	mkdir -p bin/random_weather
	GOOS=linux GOARCH=amd64 go build -o bin/random_weather/bootstrap lambda_functions/random_weather/main.go
	zip -j bin/random_weather/main.zip bin/random_weather/bootstrap

.PHONY: compile-update_weather_from_csv
compile-update_weather_from_csv:
	mkdir -p bin/update_weather_from_csv
	GOOS=linux GOARCH=amd64 go build -o bin/update_weather_from_csv/bootstrap lambda_functions/update_weather_from_csv/main.go
	zip -j bin/update_weather_from_csv/main.zip bin/update_weather_from_csv/bootstrap

.PHONY: compile-get_city_weather
compile-get_city_weather:
	mkdir -p bin/get_city_weather
	GOOS=linux GOARCH=amd64 go build -o bin/get_city_weather/bootstrap lambda_functions/get_city_weather/main.go
	zip -j bin/get_city_weather/main.zip bin/get_city_weather/bootstrap

.PHONY: compile-put_city_weather
compile-put_city_weather:
	mkdir -p bin/put_city_weather
	GOOS=linux GOARCH=amd64 go build -o bin/put_city_weather/bootstrap lambda_functions/put_city_weather/main.go
	zip -j bin/put_city_weather/main.zip bin/put_city_weather/bootstrap