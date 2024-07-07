# Build stage
FROM golang:1.22.4-alpine3.20 AS builder
# create app folder
WORKDIR /app
# copy this project folders, files to the new created docker's /app directory
COPY . .
# run command. This command will create a binary project thanks to copied files
RUN go build -o main main.go
# alpine linux hasn't curl  
# add curl
#RUN apk add curl 
# use curl to download golang migration
#RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz


# Run stage
# to reduce docker project size
# it uses the builded binary file only. remove golang and all packages that are required by our project

# we use from instruction to specify the base image for this run stage
# and there is a alpine linux image here
FROM alpine:3.20
WORKDIR /app
# the copy file in the builder stage so we specify with --from keyword.
# builder stage has a binary golang app in /app/main. We take it and copy to current stage /app folder
COPY --from=builder /app/main .
# put migrate file same workdir folder 
#COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

# it informs the docker the container listens on the 8080 port at runtime
EXPOSE 8080
# it starts golang app
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh"]