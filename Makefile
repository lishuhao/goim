# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

#all: test build
build:
	rm -rf target/
	mkdir target/
	cp cmd/comet/comet-example.toml target/comet.toml
	cp cmd/logic/logic-example.toml target/logic.toml
	cp cmd/job/job-example.toml target/job.toml
	$(GOBUILD) -o target/comet cmd/comet/main.go
	$(GOBUILD) -o target/logic cmd/logic/main.go
	$(GOBUILD) -o target/job cmd/job/main.go
#
#test:
#	$(GOTEST) -v ./...
#
#clean:
#	rm -rf target/
#
run:
	nohup target/logic -alsologtostderr -conf=target/logic.toml -region=sh -zone=sh001 -deploy.env=dev -weight=10 2>&1 > target/logic.log &
	nohup target/comet -alsologtostderr -conf=target/comet.toml -region=sh -zone=sh001 -deploy.env=dev -weight=10 -addrs=127.0.0.1 -debug=true 2>&1 > target/comet.log &
	nohup target/job -alsologtostderr -conf=target/job.toml -region=sh -zone=sh001 -deploy.env=dev 2>&1 > target/job.log &

stop:
	pkill -f target/logic
	pkill -f target/job
	pkill -f target/comet

#----------------------------

up: start_discovery build  run
down: stop stop_discovery
restartComet:
	pkill -f target/comet
	$(GOBUILD) -o target/comet cmd/comet/main.go
	nohup target/comet -alsologtostderr -conf=target/comet.toml -region=sh -zone=sh001 -deploy.env=dev -weight=10 -addrs=127.0.0.1 -debug=true 2>&1 > target/comet.log &

start_discovery:
	nohup dist/discovery/discovery111 -conf=dist/discovery/discovery-example.toml 2>&1 > dist/discovery/discovery.log &
stop_discovery:
	pkill -f dist/discovery/discovery111

start_dc:
	cd compose && docker-compose up -d
stop_dc:
	cd compose && docker-compose down