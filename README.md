# prime_number

## How to use
- Download/clone the project then import Prime_Number_Service.postman_collection.json file to your Postman Collection
- Edit environment variable in postman: set var SERVER into ec2-13-229-81-27.ap-southeast-1.compute.amazonaws.com:8000 or 13.229.81.27:8000
- You can try to call this API: {{SERVER}}/v1/health 
- If you can see response like this, server is working properly: {"status": "ok"}
- Call this API to get token, almost every APIs require token to work: {{SERVER}}/v1/users/token
- APIs explain:
    - POST | {{SERVER}}/v1/users | create a new user
    - GET | {{SERVER}}/v1/users | get all users
    - GET | {{SERVER}}/v1/users/:id | get specific user
    - PUT | {{SERVER}}/v1/users/:id | edit user's info
    - DEL | {{SERVER}}/v1/users/:id | delete user (only admin can do this action)
    - POST | {{SERVER}}/v1/requests | create a prime number request (edit any number you want in the field e.g. {"send_number": 123})
    - GET | {{SERVER}}/v1/requests | get all requests
    - GET | {{SERVER}}/v1/requests/:id | get all requests by a user

- About profiling, you can go to http://13.229.81.27:6060/debug/pprof to take a look at server's profiles
- About metrics, you can go to http://13.229.81.27:6060/debug/vars to see the debug server
- About tracing, you can go to http://13.229.81.27:9411/zipkin to observe all the traced requests


### Schema

<img src="./Prime number .png">


### Middlewares diagram

<img src="./middlewares.png">

### TODO List

- [x] Find a solution to find the nearest prime number
- [x] Choose tech stacks. Backend: Golang, DB: Postgres, Git: Github, Server: AWS EC2, CI/CD: Github Actions
- [x] Design DB Schema and generate SQL queries
- [x] Set up DB run on a Docker container
- [x] Make DB migrations
- [x] Create a Makefile to shorten the commands
- [x] Create a simple server
- [x] Handle user CRUD, intergrate find the nearest prime number function, write tests
- [x] implement health check, profilling, add middleware, metrics, request logging
- [x] Implement Authentication and Authorization
- [x] Add tracing
- [x] Create a Dockerfile, and docker-compose 
- [x] Deploy project on an AWS EC2 manually
- [ ] Set up Github Actions
