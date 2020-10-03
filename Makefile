CMD := todosvr
GOARCH := amd64
GOOS := darwin

build: swag-init
	go build -o ./target/${GOARCH}_${GOOS}/${CMD} ./

swag-init:
	swag init
