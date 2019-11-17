###Test task

Go version used: `1.13.4`

This is a test solution that caches data from http-data source to redis.
I tried to consider most corner cases in this solution and created additional interfaces/functions to create a thin wrap
over existing and standard implementations for me to able to test business logic.

###Solution peculiarities
* I tried not to unmarshal/marshal structures that are being passed to external API, as this service shouldn't know about
implementation details of external API. This will allow us not to modify our code, when implementation of external API changes,
for example, a new field is added to response.
* As a result of previous point, this service only knows about `contact_id` that is used as a key for values in Redis.
* It doesn't cache headers, received from `GET` response. Task didn't mention whether it should be 
done or not, so I decided not to cache them. That's why responses might differ - cache-miss `GET` requests will have all
headers, provided by external API, cache-hit `GET` requests will have only `Content-Type` header, besides default ones.
If it is required to cache headers - can be implemented rather easily, as there are corresponding abstractions in the code.

### Contents
This project contains packages that can be treated as utility packages:
* [logs](logs/README.md) - logging interface and implementations
* [utils](utils/README.md) - utility functions
* [mocks](mocks/README.md) - mocks for unit tests
* consts - list of consts used in this repo

and solution package:
* [middleware](middleware/README.md) - core interfaces and implementations to solve the task

###Unit tests
Package `middleware/sources` is covered with tests, as it contains a core business logic.

Package `miggleware/handles` contains unit-tests only for a common part of handlers.

###Things to improve:
* add more unit-tests and reduce code duplication in existing tests. It is possible to add tests in `middleware` package
and to cover code with tests in `utils` and `logs` packages.
* add more logging. To keep code simple, I did less logging.
* add more options to redis config.

###Things to keep in mind:
* if it is a production release, one should consider using https and providing `OPTIONS`-method for existing handlers.
 This can be done either using a proxy/load-balancer before hitting this service or (worst case scenario) via using
 ambassador template with `nginx` inside container, that will run in the same network-namespace as a container with this service.


##How to run the service:
Modify file `./config.json`:
* `api_url` - URL to external API (`https://my.test.com/v1/api/entity`).
* `cache_ttl_seconds` - for how long we should keep cached value (in seconds).
* `redis` - block of redis configuration. Supports only address (`host:port`) and DB. **Redis password is provided via
command line**.
* `bind` - which IP and port should be used by the service.
* `app_timeout_seconds` - when `SIGINT` or `SIGTERM` is caught, application is informed and should stop withing this
time interval, otherwise it will be killed.

###Source code:
`go build && ./test -config=./config.json -redispwd='securepassword'`
###Docker-way:
We will pull redis container and run it without any authentication.
1. Create a redis container:
    * `docker run -d --name redistest redis:latest`
2. Build container with service inside (from root of this repo):
    * `docker build . -t cachetest`
3. Run container with service:
    * `docker run -d --network=container:redistest --name cachetest cachetest -redispwd=""`

Those commands will spin a redis container without authentication, will create a container with the service,
building it from source code and will spin that container with the service inside, providing same network namespace as
redis container will have.

To remove everything execute following commands:
1. `docker rm -f cachetest && docker rmi cachetest`
2. `docker rm -f redistest && docker rmi redis:latest`

If you have redis wuth authentication running in another container or on host, you can execute `docker run` for `cachetest`,
providing redis password in `-redispwd` parameter:

`docker run -d  --name cachetest cachetest -redispwd="securepassword"`

Also you can use precompiled container ([this one](https://hub.docker.com/repository/docker/coldze/svctest)), substituting `config.json` file:

`docker run -d -v $(pwd)/build/config.json:/go/src/app/config.json --network=container:redistest --name cachetest coldze/svctest -redispwd=""`

###Sample tests:
Those commands can be executed both from host and from inside container, but from root of repo (json files are required for post/put methods):
* GET: `curl http://<service-container-ip>/v1/contact/<contact-id>`
* POST: `curl -X "POST" -H "Content-Type: application/json" -d @test_data/post.json http://<service-container-ip>/v1/contact`
* PUT: `curl -X "POST" -H "Content-Type: application/json" -d @test_data/put.json http://<service-container-ip>/v1/contact`

##How to run unit-tests:
From root of repo run following command:

`go test ./...`