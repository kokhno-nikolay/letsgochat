<h1 align="center">
  <img src="https://raw.githubusercontent.com/create-go-app/cli/master/.github/images/cgapp_logo%402x.png" width="224px"/><br/>
  Let's Go Chat
</h1>

<p align="center"><a href="https://pkg.go.dev/github.com/create-go-app/cli/v3?tab=doc" target="_blank"><img src="https://img.shields.io/badge/Go-1.17+-00ADD8?style=for-the-badge&logo=go" alt="go version" /></a>&nbsp;<img src="https://img.shields.io/badge/license-apache_2.0-red?style=for-the-badge&logo=none" alt="license" />&nbsp;<a href="https://pkg.go.dev/github.com/kokhno-nikolay/lets-go-chat"><img src="https://pkg.go.dev/badge/github.com/kokhno-nikolay/lets-go-chat.svg" alt="Go Reference"></a></p>

## ⚡ Quick start
First of all, [download](https://golang.org/dl/) and install **Go**. Version `1.17` or higher is required.

## ⚙ Config (.env)
Create <b>.env</b> file and add the value as in example.

| Name                   | Type | Example value |
|------------------------| ------ |  ------ |
| SECRET_KEY | string | secret |
| POSTGRES_HOST_NAME     | string | postgres |
| POSTGRES_PORT          | int | 5432 |
| POSTGRES_USERNAME      | string | postgres |
| POSTGRES_PASSWORD      | string | postgres |
| POSTGRES_DATABASE_NAME | string | lets_go_chat |

## 🧩 How to run
1) Clone project
```
git clone https://github.com/kokhno-nikolay/letsgochat.git
```
2) Start docker-compose
```
go run main.go
```

##  ‍🚀 API
<b>URL</b> - baseURL/v1/
1) Registration
```
POST https://letsgochat.herokuapp.com/user -H 'Content-Type: application/json' -d '{"username":"someusername","password":"random-password"}'
```
2) Login
```
POST https://letsgochat.herokuapp.com/user/login -H 'Content-Type: application/json' -d '{"username":"someusername","password":"random-password"}'
```
3) Active
```
GET https://letsgochat.herokuapp.com/user/active -H 'Content-Type: application/json
```
4) Websocket chat
```
ws://letsgochat.herokuapp.com/chat?token=<token>
```