# Golang Backend Project

### Problem Statement

Build a to-do application with Golang as the backend
The frontend can either be CLI commands or a simple frontend interface
The todos entered in the frontend should be stored in a database (Postgres)
The data fetches from the database should be cached with Redis.
Write a basic unit test in Golang to test your connection to Postgres and Redis
README to run this setup locally

<br>

### My Solution

#### Steps to run

1. Clone the repo
2. Run the command **`docker-compose up`** in the folder having the docker-compose.yaml file

(If you see the following *Postgres DB init automigration completed* and *Redis cache init was completed*, then you are good to go.) <br>
The app listens on port **8081**.
<br>


#### Explanation of project

I started the project with the 3 basic files- main.go, go.mod and go.sum. I'm making use of gorilla mux for the multiplexer, and the code for the handler functions can be found in the controllers package. <br>
I started out with first writing the golang app with data being persisted only in the postgres db, and then eventually added redis as the cache layer. So for requests such as getTodo() which would earlier directly contact the postgresdb would now first check for the entry in the redis cache. We store and set elements in the cache when the user first creates a new todo, and then if they update a todo task, and also when user makes a get request and if it wasn't found in the cache. <br>
The basic layout of the 3 main components in the projecr are as follows: 
<br>
![architecture](https://github.com/geegatomar/Golang-Redis-Postgres-Project/blob/master/images/architecture.png?raw=true)

I then dockerized the application, and wrote the Dockerfile and docker-compose.yaml files respectively. <br>
The **postman collection** for the five major APIs can be found here: https://www.postman.com/collections/b38404c7f31cdd89d33b



#### Simple Testing
Since the problem statement asks us to write only very basic tests, to test the apis along with the database connections, for the sake of simplicity I have not written tests testing each layer individually and making a mock db, etc. <br>
Instead, we are testing each of the routes end-to-end here in the file *main_test.go*, and hence the server must be up and running when you execute these test commands. <br>
I have included a setup and teardown function for some of the tests which require some data to be already present in the db.
Command to test:  `go test -v`  <br>
![testing](https://github.com/geegatomar/Golang-Redis-Postgres-Project/blob/master/images/testing.png?raw=true)



#### Caveats & Future improvements
1. Testing can be improved with use of mockdb, and we can also perform benchmarking.
2. For scalability, We can spin up multiple containers and balance the load using nginx (like I had done here for another project: https://github.com/geegatomar/Load-Balancing-for-Image-processing-systems).  <br>
We can also add volumes in docker to persists the storage.
3. Efficient handling of bad requests, and handling edge cases associated with the same.
