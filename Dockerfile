ARG GO_VERSION=1.19.1 
FROM public.ecr.aws/docker/library/golang:${GO_VERSION}-alpine AS build

WORKDIR /go/src/helloworld

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build \
   -v \
   -o /helloworld

FROM scratch
 
COPY --from=build /helloworld /helloworld
 
ENTRYPOINT ["/helloworld"]
