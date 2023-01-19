FROM golang:1.19 as builder

WORKDIR /workspace
COPY go.* .
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /docmerge cmd/docmerge/main.go


# Use distroless as minimal base image to package the project
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /docmerge .
ENTRYPOINT ["/docmerge"]