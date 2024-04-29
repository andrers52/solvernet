FROM golang:latest 

LABEL base.name="solvernet"

WORKDIR /app 

COPY . . 

RUN go build . 

EXPOSE 3002 

ENTRYPOINT [ "./blockchain" ]
