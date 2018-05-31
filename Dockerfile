FROM golang:1.10-alpine3.7 AS build
RUN apk add --no-cache make
WORKDIR /go/src/github.com/sgreben/versions/
COPY . .
ENV CGO_ENABLED=0
RUN make binaries/linux_x86_64/versions

FROM scratch
COPY --from=build /go/src/github.com/sgreben/versions/binaries/linux_x86_64/versions /versions
ENTRYPOINT [ "/versions" ]
