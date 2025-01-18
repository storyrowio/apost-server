FROM golang:1.23-alpine

WORKDIR /app
COPY . .
RUN go mod tidy

COPY *.go ./

RUN rm -f .env .env.* .env-*
RUN rm -rf .git .gitignore

RUN go build -o /apost

EXPOSE 8000

CMD [ "/apost" ]