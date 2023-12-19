# syntax=docker/dockerfile:1
FROM node:19 AS build
WORKDIR /app
COPY ["client/package.json", "client/package-lock.json*", "./"]
RUN npm i
COPY client/ .
# https://github.com/parcel-bundler/parcel/issues/7126
RUN rm -rf .parcel-cache/
RUN npm run build

FROM golang:1.21-alpine
COPY --from=build /app/dist /app/client/dist
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY main.go ./
RUN go build -o /server
CMD [ "/server" ]
