FROM golang:alpine AS build

RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=on

ADD playground playground
ADD webapi webapi
ADD main.go main.go
ADD go.mod go.mod
ADD go.sum go.sum
ADD examples examples

RUN go mod vendor
RUN go mod verify
RUN go mod tidy
RUN GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -tags ptr64 -o /build/cxplayground /build/main.go

FROM busybox

WORKDIR /var/cxplayground

COPY --from=build /build/examples /var/cxplayground/examples
COPY --from=build /build/cxplayground /var/cxplayground/cxplayground
COPY dist /var/cxplayground/dist

EXPOSE 5336

ENTRYPOINT ["/var/cxplayground/cxplayground"]
