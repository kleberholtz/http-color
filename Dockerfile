FROM golang:alpine AS build

WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /server .

FROM scratch
COPY --from=build /server /server
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/server"]
