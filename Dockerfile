# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:alpine AS build

WORKDIR /src
ARG TARGETOS
ARG TARGETARCH
COPY go.mod main.go ./
RUN --mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
	go build -trimpath -buildvcs=false -ldflags="-s -w" -o /server .

FROM scratch
COPY --from=build /server /server
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/server"]
