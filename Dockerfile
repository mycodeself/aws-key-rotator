FROM golang:1.15 AS builder

WORKDIR /src

COPY . .

RUN make build

FROM scratch AS app

COPY --from=builder /src/bin/aws-key-rotator /usr/bin/

CMD [ "aws-key-rotator" ]