PROJ=fake-airline-info-api

export APP_NAME=$(PROJ)
export DEBUG=true
export CONFIG_DOCKER=config
export CONFIG_LOCAL=config-local

APP_PORT=3000
APP_NAME_ON_LIARA=mock-flight-booking
deploy_liara_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiI2MzFiMWVmZTFlYTQ1ZTgyMjljZmJlOTEiLCJpYXQiOjE2ODkzMzEyNTJ9.ikHUXh5e4-DdazjtNXJLoZsf8gLPal4odWMlltNFG48

export APP_CONFIG=$(CONFIG_LOCAL)
export APP_ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

build:
	go build -race -o ./$(APP_NAME) ./main.go

run:
	go run -race ./main.go

test:
	go test -v -race ./... -count=1

format:
	find . -name '*.go' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

fmt: format

clean:
	rm -f ./$(APP_NAME)

docker:
	export APP_CONFIG=$(CONFIG_DOCKER)
	docker-compose up --build

deploy:
	liara deploy --app=$(APP_NAME_ON_LIARA) --region=iran \
	--api-token=$(deploy_liara_token) --port=$(APP_PORT)
