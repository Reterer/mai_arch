##
## STEP 1 - BUILD
##

# specify the base image to  be used for the application, alpine or ubuntu
FROM golang:1.21-alpine AS build

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY . .


# download Go modules and dependencies
RUN go mod download
# compile application
RUN go build -o app client_service/main.go

##
## STEP 2 - DEPLOY
##
FROM alpine:3.18
EXPOSE 8080 
WORKDIR /app

COPY client_service/config/config.json config/
COPY --from=build ./app .

ENTRYPOINT [ "./app" ]

