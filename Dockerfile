FROM golang:latest as stage

WORKDIR /build

COPY . .

RUN go build -o app

FROM golang:latest

COPY --from=stage /build/app /go/bin/app

ENV forum=
ENV interval=
ENV GPT_SECRET=
ENV ACCESS_ID_AWS=
ENV SECRET_ID_AWS=
ENV MONGO_URI=
ENV DATABASE_NAME=Posts

ENTRYPOINT [ "/go/bin/app" ]