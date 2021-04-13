NAME := gc-gsm-secrets
build: clean
	go build -o build/$(NAME) main.go

clean:
	rm -rf build/