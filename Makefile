VERSION?=dev
GOOS?=linux
BUILD_PATH?=./target
FORMATS_PATH?=./target/internal/formats/data_getter/

build:
	rm -f $(BUILD_PATH)/format-bot
	go build $(GO_TAGS) $(GO_LDFLAGS) -o $(BUILD_PATH)/format-bot ./
	cp -v config-dev.yml $(BUILD_PATH)/
	mkdir -p $(FORMATS_PATH)/
	cp -v ./internal/formats/data_getter/questions.json $(FORMATS_PATH)

build-dev:
	rm -f $(BUILD_PATH)/format-bot
	go build $(GO_TAGS) $(GO_LDFLAGS) -o $(BUILD_PATH)/format-bot ./
	cp -v config-dev.yml $(BUILD_PATH)/
	mkdir -p $(FORMATS_PATH)/
	cp -v ./internal/formats/data_getter/questions.json $(FORMATS_PATH)

clean:
	rm -rf $(BUILD_PATH)
