VERSION?=dev
GOOS?=linux
BUILD_PATH?=./target

build:
	rm -f $(BUILD_PATH)/format-bot
	go build $(GO_TAGS) $(GO_LDFLAGS) -o $(BUILD_PATH)/format-bot ./
	cp -v config.yml $(BUILD_PATH)/
	mkdir -p $(BUILD_PATH)/formats/
	cp -rv formats/formats.json $(BUILD_PATH)/formats/

clean:
	rm -rf $(BUILD_PATH)
