
export GOPROXY=goproxy.cn,direct
GOFLAGS=""

OBJ = go-webrdp-demo

all: $(OBJ)

$(OBJ):
	go build -ldflags $(GOFLAGS) -o $(OBJ) main.go

clean:
	rm -fr $(OBJ)

-include .deps

dep:
	/bin/echo -n "$(OBJ):" > .deps
	find . -path ./vendor -prune -o -name '*.go' -print | awk '{print $$0 " \\"}' >> .deps
	/bin/echo "" >> .deps
