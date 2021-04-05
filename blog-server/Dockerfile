FROM golang:1.15-alpine3.12 AS builder

LABEL stage="builder"

ENV GO111MODULE=on \
    CGO_ENABLED=1

WORKDIR /blogserver

COPY . .

RUN apk add --no-cache gcc libc-dev \ 
    && go mod download

RUN CGO_ENABLED=0 GOOS=linux go build .     

FROM alpine AS final

LABEL stage="final"

WORKDIR /blogapp

COPY --from=builder /blogserver/blog .

EXPOSE 50051

CMD [ "./blog" ]