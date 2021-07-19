# Go Realtime Chat

Simple realtime chat server built using golang websockets.

## Run Locally

Clone the project

```bash
  git clone https://github.com/jwambugu/go-realtime-chat.git
```

Go to the project directory

```bash
  cd my-project
```

Start the webserver

```go
  go run cmd/web/*
```

## Chat Application

To access the chat application, make sure you start the webserver then open ```http://localhot:8080```

![Chatroom](https://github.com/jwambugu/go-realtime-chat/blob/main/screenshots/Screenshot%202021-07-19%20at%2008.53.32.png?raw=true)

Enter the chatroom name then create the room then have fun chatting 🎉

![Chat](https://github.com/jwambugu/go-realtime-chat/blob/main/screenshots/output.gif?raw=true)

## Todo

- Use the username in the chat messages.

- Store the messages in the DB.

    

