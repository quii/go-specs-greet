# go-specs-greet

Source code for the chapter (currently WIP) "Scaling Acceptance Tests (and light intro to gRPC)"

**Idea**: Dog-food this chapter on twitch before releasing.

# Learn Go with Tests - Scaling Acceptance Tests (and light intro to gRPC)

This is a follow up to [Intro to acceptance tests](https://quii.gitbook.io/learn-go-with-tests/testing-fundamentals/intro-to-acceptance-tests)

When written well, acceptance tests are an important part of a systems test suite. They can be used at different abstraction layers to give yourself confidence that your system works how you need it to.



## Ideas / things left to write about

- Discuss Dave Farleys acceptance test youtube vid, and reference gopherconuk talk by Riya and I
- Don't write an acceptance test for everything, reference the test pyramid again
  - Adding language switch would demand a change in the spec as we're changing the API of the system
  - Subsequent languages should be done via unit tests




## Things reader will learn

- How to use specifications / drivers to decouple the accidental and essential complexity
  - Normally when you're solving someone's problem, you're dealing with essential complexity, try and express that in the specification
- Top-down [GOOS](http://www.growing-object-oriented-software.com)-thinking
  - Start with hello, world, build from there
- Intro to gRPC

## Prerequisite material

There's lots of ideas and inspiration for this chapter, a lot of it born from many years of frustration with acceptance tests causing lots of issues! The main two videos I would recommend you watch are

- Dave Farley - Acceptance Tests
- Nat Pryce - E2E functional tests that can run in milliseconds
- GOOS - Nat Pryce & Steve Freeman



## Anatomy of bad acceptance tests

For many years, I've worked for several companies and teams. Each of them recognised the need for acceptance tests, some way to test a system from a user's point of view and verify it works how it's intended, but almost without exception, the cost of these tests become a real problem for the team.

- Slow
- Still have numerous bugs
- Brittle, expensive to maintain, seem to make changing the software harder than it aught to be
- Can only run in a very specific environment, causing poor feedback loops


## Anatomy of good acceptance tests

TODO:

- Talk about separation of domain and accidental complexity
- Importance of polymorphism and re-use. Talk about how if the specs represent the "truths" of how we want the system to behave, we can verify it at various abstractions, from our domain code, to the application as a black box, to it being deployed to a staging environment, even to running it against live where CDNs can have an effect (bad cache headers or whatever)



## Let's go

Create a new project

`go mod init github.com/quii/go-specs-greet`

Make a folder `specifications` to hold our specification, and add a file `greet.go`

```go
package specifications

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

type Greeter interface {
	Greet() (string, error)
}

func GreetSpecification(t testing.TB, greeter Greeter) {
	got, err := greeter.Greet()
	assert.NoError(t, err)
	assert.Equal(t, got, "Hello, world")
}
```

My IDE (Goland) takes care of the fuss of adding dependencies for me, but if you need to do it manually you'd do

`go get github.com/alecthomas/assert/v2`

Given Farley's acceptance test design, we now have a specification which is decoupled from implementation. It doesn't know, or care about _how_ we `Greet`, it's just concerned with the logic. This "logic" isn't much right now, but we'll expand upon the spec to add more functionality as we further iterate.

At this point, this level of ceremony to decouple our specification from implementation might make some people accuse us of "overly abstracting"; I promise you that acceptance tests that are too coupled to implementation become a real burden on engineering teams. I am confident to assert that most acceptance tests out in the wild are expensive to maintain, due to this inappropriate coupling; rather than the reverse, of being overly abstract.

We can use this specification to verify any "system" that can `Greet`.

### First system: HTTP API

Our requirement is to provider a greeter service over HTTP. So we'll need to create:

1. A **driver**. In this case, the way one works with a HTTP system is using a **HTTP client**. This code will know how to work with our API. Drivers implement the interface that specifications define.
2. A HTTP server with a greet API
3. A test, which is responsible for managing the life-cycle of spinning up the server, and then plugging the driver into the specification to run it as a test

## Write the test first

The initial process for creating a black-box test that compiles and runs your program, executes the test and then cleans everything up can be quite labour intensive. That's why it's preferable to do it at the start of your project on a very small amount of functionality. I typically start all my projects with a "hello world" server implementation, with all of my tests set up, ready for me to build the real functionality easily.

Most development teams these days are shipping using docker, so our acceptance tests will test a docker image we'll build of our program.

To help us use Docker in our tests, we're going to use [Testcontainers](https://golang.testcontainers.org).

`go get github.com/testcontainers/testcontainers-go`

Create some structure to house our program we intend to ship

`mkdir -p cmd/http_server`

Inside the new folder, create a new file and add the following

`greeter_http_server_test.go`

```go
package main_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	go_specs_greet "github.com/quii/go-specs-greet"
	"github.com/quii/go-specs-greet/specifications"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGreeterServer(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../.",
			Dockerfile: "./cmd/http_server/Dockerfile",
		},
    ExposedPorts: []string{"8080:8080"},
		WaitingFor:   wait.ForHTTP("/").WithPort("8080"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	driver := go_specs_greet.Driver{BaseURL: "http://localhost:8080"}
	specifications.GreetSpecification(t, driver)
}
```

Notes:

- Most of the code is dedicated to building the Docker image of our web server and then launching a container from it
- We're going to allow our driver to be configurable with the `BaseURL` field. This'll allow us to re-use the driver in different environments, such as staging, or even production.

## Try to run the test

```
./greeter_server_test.go:46:12: undefined: go_specs_greet.Driver
```

We're still practicing TDD here! It's a big first step we have to make, we need to make a few files and write maybe more code than we're typically used to, but when you're first starting this is often the case. It's so important we try and remember the rules of the red step.

> Commit as many sins as neccessary to get the test passing

## Write the minimal amount of code for the test to run and check the failing test output

Hold your nose, and remember we can refactor when the test is passing. Here's the code for our driver in `driver.go`

```go
package go_specs_greet

import (
	"io"
	"net/http"
)

type Driver struct {
	BaseURL string
}

func (d Driver) Greet() (string, error) {
	res, err := http.Get(d.BaseURL + "/greet")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	greeting, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(greeting), nil
}
```

Notes:

- You could argue that perhaps I should be writing tests to drive out the various `if err != nil`, but in my experience so long as you're not doing anything with the `err`, tests that say "you return the error you get" are fairly low value.
- **You shouldn't use the default HTTP client**. Later we'll pass in a HTTP client so it can be configured with timeouts e.t.c., but for now we're just trying to get ourselves to a passing test

Try and run the tests again, they should now compile, but not pass.

```
=== RUN   TestGreeterHandler
2022/09/10 18:49:44 Starting container id: 03e8588a1be4 image: docker.io/testcontainers/ryuk:0.3.3
2022/09/10 18:49:45 Waiting for container id 03e8588a1be4 image: docker.io/testcontainers/ryuk:0.3.3
2022/09/10 18:49:45 Container is ready id: 03e8588a1be4 image: docker.io/testcontainers/ryuk:0.3.3
    greeter_server_test.go:32: Did not expect an error but got:
        Error response from daemon: Cannot locate specified Dockerfile: ./cmd/http_server/Dockerfile: failed to create container
--- FAIL: TestGreeterHandler (0.59s)

```

We need to create a Dockerfile for our program. Inside our `http_server` folder, create a `Dockerfile` and add the following

```dockerfile
FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o svr cmd/http_server/*.go

EXPOSE 8080
CMD [ "./svr" ]
```

Don't worry too much about the details here, it can be refined and optimised, but for this example, it'll suffice. The advantage of our approach here is we can later improve our Dockerfile and have a test to prove it works as we intend it to. This is the real strength of having black-box tests!

Try and run the test again and it should complain about not being able to build the image. That's because we haven't added a program yet!

For the test to fully execute, we'll need to create a program that listens on `8080`, but **that's all**. Stick to the TDD discipline, don't write the production code that would make the test pass until we've verified the test fails as we'd expect.

Create a `main.go` inside our `http_server` folder with the following

```go
func main() {
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	})
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
```

```
    greet.go:16: Expected values to be equal:
        +Hello, World
        \ No newline at end of file
--- FAIL: TestGreeterHandler (2.09s)
```

## Write enough code to make it pass

Update the handler to behave how our specification wants it to

```go
func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "Hello, world")
	})
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
```

## Refactor

Whilst this technically isn't a refactor, we shouldn't rely on the default HTTP client, so let's change our client so it can be supplied one; which our test will give to it.

```go
type Driver struct {
	BaseURL string
	Client *http.Client
}

func (d Driver) Greet() (string, error) {
	res, err := d.Client.Get(d.BaseURL + "/greet")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	greeting, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(greeting), nil
}
```

Update the creation of the driver to pass in a client.

```go
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	driver := go_specs_greet.Driver{BaseURL: "http://localhost:8080", Client: &client}
	specifications.GreetSpecification(t, driver)
}
```

It's good practice to keep `main.go` as simple as possible, it only really aught to be concerned with piecing together the building blocks you make in to an application.

Create a file called `handler.go` and move our code into there

```go
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world")
}
```

Update `main.go` to import and use the handler instead

```go
package main

import (
	"net/http"

	go_specs_greet "github.com/quii/go-specs-greet"
)

func main() {
	handler := http.HandlerFunc(go_specs_greet.Handler)
	http.ListenAndServe(":8080", handler)
}
```

## Reflect

The first step felt like an effort. We've made a number of `go` files to create and test a HTTP handler that returns a hard-coded string. This "iteration 0" ceremony and setup though will serve us well for further iterations.

Adding or changing functionality should be simple, and controlled by driving it through the specification and dealing with whatever changes it drives us to do. Now the `DockerFile` and `testcontainers` are set up for our acceptance test, we shouldn't have to change these files unless the way we construct our application changes.

We'll see this with our next requirement, greet a particular person.

## Write the test first

Edit our specification

```go
package specifications

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

type Greeter interface {
	Greet(name string) (string, error)
}

func GreetSpecification(t testing.TB, greeter Greeter) {
	got, err := greeter.Greet("Mike")
	assert.NoError(t, err)
	assert.Equal(t, got, "Hello, Mike")
}

```

To allow us to greet specific people, we need to change the interface to our system to accept a `name` parameter.

## Try to run the test

```
./greeter_server_test.go:48:39: cannot use driver (variable of type go_specs_greet.Driver) as type specifications.Greeter in argument to specifications.GreetSpecification:
	go_specs_greet.Driver does not implement specifications.Greeter (wrong type for Greet method)
		have Greet() (string, error)
		want Greet(name string) (string, error)
```

The change in the specification has meant our driver needs to be updated.

## Write the minimal amount of code for the test to run and check the failing test output

```go
func (d Driver) Greet(name string) (string, error) {
	res, err := d.Client.Get(d.BaseURL + "/greet?name=" + name)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	greeting, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(greeting), nil
}
```

The test should now run

```
    greet.go:16: Expected values to be equal:
        -Hello, world
        \ No newline at end of file
        +Hello, Mike
        \ No newline at end of file
--- FAIL: TestGreeterHandler (1.92s)
```

## Write enough code to make it pass

```go
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s", r.URL.Query().Get("name"))
}
```

## Refactor

In [HTTP Handlers Revisited](https://github.com/quii/learn-go-with-tests/blob/main/http-handlers-revisited.md) we discussed how important it is for HTTP handlers should only be response for handling HTTP concerns, any kind of "domain logic" should live outside of the handler. This allow us to develop domain logic in isolation of HTTP, making it simpler to test and understand.

Let's pull apart these concerns.

```go
func Handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprint(w, Greet(name))
}
```

In `greet.go`

```go
func Greet(name string) string {
	return fmt.Sprintf("Hello, %s", name)
}
```

## A slight diversion in to the "adapter" design pattern

Now that we've separated our domain logic of greeting people into a separate function, we are now free to write unit tests for our greet function; certainly a lot simpler than testing it through a specification, that goes through a driver, that hits a web server, to finally get a string!

Wouldn't it be nice if we could re-use our specification here too. After-all, the point of the specification is it's decoupled from implementation details.

Let's give it a go in `greet_test.go`

```go
func TestGreet(t *testing.T) {
	specifications.GreetSpecification(t, go_specs_greet.Greet)
}
```

This would be nice, but it doesn't work

```
./greet_test.go:11:39: cannot use go_specs_greet.Greet (value of type func(name string) string) as type specifications.Greeter in argument to specifications.GreetSpecification:
	func(name string) string does not implement specifications.Greeter (missing Greet method)
```

Our specification wants something that has a method `Greet()` not a function.

This is frustrating, we have a thing that we "know" is a `Greeter`, but it's not quite in the right **shape** for the compiler to let us use it. This is what the **adapter** pattern caters for.

> In [software engineering](https://en.wikipedia.org/wiki/Software_engineering), the **adapter pattern** is a [software design pattern](https://en.wikipedia.org/wiki/Software_design_pattern) (also known as [wrapper](https://en.wikipedia.org/wiki/Wrapper_function), an alternative naming shared with the [decorator pattern](https://en.wikipedia.org/wiki/Decorator_pattern)) that allows the [interface](https://en.wikipedia.org/wiki/Interface_(computer_science)) of an existing [class](https://en.wikipedia.org/wiki/Class_(computer_science)) to be used as another interface.[[1\]](https://en.wikipedia.org/wiki/Adapter_pattern#cite_note-HeadFirst-1) It is often used to make existing classes work with others without modifying their [source code](https://en.wikipedia.org/wiki/Source_code).

This is a lot of fancy words, for something that is quite simple. Which is often the case with design patterns, which is why people tend to roll their eyes at them. The value of design patterns is not specific implementations, but a language to describe certain solutions to common problems engineers face. If you have a team that has a shared vocabulary, it reduces the friction in communication.

Add this code in `greet.go`

```go
type GreetAdapter func(name string) string

func (g GreetAdapter) Greet(name string) (string, error) {
	return g(name), nil
}
```

We can now use our adapter in our test to plug our `Greet` function into the specification.

```go
func TestGreet(t *testing.T) {
	specifications.GreetSpecification(
		t,
		gospecsgreet.GreetAdapter(gospecsgreet.Greet),
	)
}
```

## Reflect

This felt simple right? OK, maybe it was simple due to the nature of the problem, but this method of work gives you discipline, a simple repeatable way of designing your code from top to bottom.

- Analyse your problem and identify a small improvement to your system that pushes you in the right direction
- Change the spec
- Follow the compilation errors until the test runs
- Update your implementation
- Refactor

After the pain of the first iteration, we didn't have to edit our acceptance test code at all because we have the seperation of specifications, drivers and implementation. Changing our specification required us to update our driver, and finally our implementation; but the boilerplate code around _how_ to spin up the system as a contaiiner was unaffected.

Even with the overhead of building a docker image for our application, and spinning up the container, the feedback loop for testing our **entire** application is very tight:

```
quii@Chriss-MacBook-Pro go-specs-greet % go test ./...
ok  	github.com/quii/go-specs-greet	0.181s
ok  	github.com/quii/go-specs-greet/cmd/httpserver	2.221s
?   	github.com/quii/go-specs-greet/specifications	[no test files]
```

Now, imagine your CTO has now decided gRPC is _the future_. She wants you to expose this same functionality over a gRPC server, whilst maintaining the existing HTTP server.

This is an example of **accidental complexity**. Accidental complexity is the complexity we have to deal with because we're working with computers, stuff like networks, disks, APIs, e.t.c. **Essential complexity** is sometimes referred to as "domain logic", it's the inescapable rules and truths within the domain you work in. They should be expressable to a non-technical person, and it's valuable to model them in our systems both in **specifications** and **domain code, that is decoupled from accidental complexity**. Many repository structures and design patterns are mainly dealing with this concern. For instance "ports and adapters" asks that you separate out your domain code from anything to do with accidental complexity, that code lives in an "adapters" folder.

Sometimes, it makes sense to do some refactoring before making a change

> Make the change easy, then make the change.

For that reason, let's gather our `http` code into a package called `httpserver` within an `adapters` folder

```
quii@Chriss-MacBook-Pro go-specs-greet % tree
.
├── adapters
│   └── httpserver
│       ├── driver.go
│       └── handler.go
├── cmd
│   └── httpserver
│       ├── Dockerfile
│       ├── greeter_server_test.go
│       └── main.go
├── go.mod
├── go.sum
├── greet.go
├── greet_test.go
└── specifications
    └── greet.go
```

Our domain code, our **essential complexity** lives at the root of our go module, and code that will allow us to use them in "the real world" are organised in to **adapters**. The `cmd` folder is where we can compose these logical groupings into useful applications, which have black-box tests to verify it all works. Nice!

Finally, we can do a _tiny_ bit of tidying up of our acceptance test. If you consider the high-level steps of our acceptance test:

- Build _some_ docker image
- Wait for it to be listening on _some_ port
- Create _some_ driver to send messages to that port
- Plug in the driver into the specification

... you'll realise we have the same requirements for an acceptance test for the gRPC server!

The `adapters` folder seems a good a place as any, so inside a file called `docker.go` , encapsulate the first 2 steps in a function that we'll re-use next.

```go
package adapters

import (
	"context"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartDockerServer(
  ctx context.Context,
	t testing.TB,
	dockerFilePath string,
	port string,
) {
	t.Helper()
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../.",
			Dockerfile: dockerFilePath,
		},
		ExposedPorts: []string{fmt.Sprintf("%s:%s", port, port)},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(5 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})
}
```

This gives us an opportunity to clean up our acceptance test a little

```go
func TestGreeterServer(t *testing.T) {
	var (
		ctx            = context.Background()
		port           = "8080"
		dockerFilePath = "./cmd/httpserver/Dockerfile"
		baseURL        = fmt.Sprintf("http://localhost:%s", port)
		driver         = go_specs_greet.Driver{BaseURL: baseURL, Client: &http.Client{
			Timeout: 1 * time.Second,
		}}
	)

	adapters.StartDockerServer(ctx, t, dockerFilePath, port)
	specifications.GreetSpecification(t, driver)
}
```

This should make writing the _next_ test simpler.

## Write the test first

You can imagine this functionality being a new adapter in to our domain code. For that reason we:

- Shouldn't have to change the specification;
- Should be able to re-use the specification;
- Should be able to re-use the domain code.

Create a new folder `grpcserver` inside `cmd` to house our new program and the corresponding acceptance test. Inside `cmd/grpc_server/greeter_server_test.go` add an acceptance test, which, not by coincedence, but by design, looks very similar to our HTTP server test.

```go
package main_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/quii/go-specs-greet/adapters"
	"github.com/quii/go-specs-greet/adapters/grpcserver"
	"github.com/quii/go-specs-greet/specifications"
)

func TestGreeterServer(t *testing.T) {
	var (
		ctx            = context.Background()
		port           = "50051"
		dockerFilePath = "./cmd/grpcserver/Dockerfile"
		addr           = fmt.Sprintf("localhost:%s", port)
		driver         = grpcserver.Driver{Addr: addr}
	)

	adapters.StartDockerServer(ctx, t, dockerFilePath, port)
	specifications.GreetSpecification(t, &driver)
}
```

The only differences are:

- We use a different docker file, because we're building a different program
- We use a different driver to plug in to the specification

## Try to run the test

```
./greeter_server_test.go:26:12: undefined: grpcserver
```

We haven't created a Driver yet, so it won't compile.

## Write the minimal amount of code for the test to run and check the failing test output

Create a `grpcserver` folder inside `adapters` and inside it create `driver.go`

```go
package grpcserver

type Driver struct {
	Addr string
}

func (d Driver) Greet(name string) (string, error) {
	return "", nil
}
```

If you run again, it should now _compile_ but not pass, because we haven't created a Dockerfile and corresponding program for it to run against.

Create a new `Dockerfile` inside `cmd/grpcserver`.

```go
FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o svr cmd/grpcserver/*.go

EXPOSE 8080
CMD [ "./svr" ]
```

And a `main.go`

```go
package main

import "fmt"

func main() {
	fmt.Println("implement me")
}
```

You should find now that the test fails because our server is not listening on the port. Now is the time to start building our client and server with gRPC

## Write enough code to make it pass

### gRPC

If you're unfamiliar with gRPC, I'd start by looking at the [gRPC website](https://grpc.io) but for the purposes of this chapter, it's just another kind of adapter in to our system, a way of other systems being able to call (**r**emote **p**rocedure **c**all) our amazing domain code.

The twist is you define a "service definition" using Protocol Buffers. You can then generate server and client code from the definition. This not only works for Go, but for most mainstream languages too. This means you can share a definition with other teams in your company who may not even write Go, and still be able to do service to service communucation very smoothly.

If you haven't used gRPC before you'll need to install a **Protocol buffer compiler** and some **Go plugins** for it. [The gRPC website has clear instructions as to how to do this](https://grpc.io/docs/languages/go/quickstart/).

Inside the same folder as our new driver, add a `greet.proto` file with the following

```protobuf
syntax = "proto3";

option go_package = "github.com/quii/adapters/grpcserver";

package grpcserver;

service Greeter {
  rpc Greet (GreetRequest) returns (GreetReply) {}
}

message GreetRequest {
  string name = 1;
}

message GreetReply {
  string message = 1;
}
```

You don't need to be an expert in Protocol Buffers to follow this definition. We're defining a service, which has a `Greet` method, and then describing the incoming and outgoing message types.

Inside `adapters/grpcserver` run the following to generate the client and server code

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    greet.proto
```

If it worked, we will have some code generated for us to use. Let's start by using the generated client code inside our `Driver`.

```go
package grpcserver

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Driver struct {
	Addr string
}

func (d Driver) Greet(name string) (string, error) {
	//todo: we shouldn't redial every time we call greet, refactor out when we're green
	conn, err := grpc.Dial(d.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := NewGreeterClient(conn)
	greeting, err := client.Greet(context.Background(), &GreetRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}

	return greeting.Message, nil
}

```

Now that we have a client, we need to update our `main.go` to create a server. Remember at this point we're just trying to get our test to pass, and not worrying about code quality.

```go
package main

import (
	"context"
	"log"
	"net"

	"github.com/quii/go-specs-greet/adapters/grpcserver"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	grpcserver.RegisterGreeterServer(s, &GreetServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

type GreetServer struct {
	grpcserver.UnimplementedGreeterServer
}

func (g GreetServer) Greet(ctx context.Context, request *grpcserver.GreetRequest) (*grpcserver.GreetReply, error) {
	return &grpcserver.GreetReply{Message: "fixme"}, nil
}
```

To create our gRPC server, we have to implement the interface it generated for us

```go
// GreeterServer is the server API for Greeter service.
// All implementations must embed UnimplementedGreeterServer
// for forward compatibility
type GreeterServer interface {
	Greet(context.Context, *GreetRequest) (*GreetReply, error)
	mustEmbedUnimplementedGreeterServer()
}
```

- Listen on the port
- We create a `GreetServer` that implements this interface, and then register it with `grpcServer.RegisterGreeterServer`, along with a `grpc.Server`.
- Use the server with the listener



It wouldn't be a huge extra effort to call our domain code inside `greetServer.Greet` rather than hard-coding `fix-me` in the message, but I'd like to run our acceptance test first, just to see everything is working end to end on a transport level

```
greet.go:16: Expected values to be equal:
-fixme
\ No newline at end of file
+Hello, Mike
\ No newline at end of file
```

Nice! We can see our driver is able to connect to our gRPC server in the test.

Now, call our domain code inside our `GreetServer`

```go
type GreetServer struct {
	grpcserver.UnimplementedGreeterServer
}

func (g GreetServer) Greet(ctx context.Context, request *grpcserver.GreetRequest) (*grpcserver.GreetReply, error) {
	return &grpcserver.GreetReply{Message: gospecsgreet.Greet(request.Name)}, nil
}
```

Finally it passes! We have an acceptance test that proves our gRPC greet server behaves how we'd like.

## Refactor

We committed a number of sins to get the test passing, but now they're passing we have the safety-net to refactor.

### Simplify main

Like before, we don't want `main` having too much code inside it, and it feels inconsistent with our other implementation. We can move our new `GreetServer` into `adapters/grpcserver` as that's definitely where it should live. In terms of cohesion if we happen to change the service definition, we want the "blast-radius" of change to be confined to that area of our code.

### Don't redial in our driver every time

Currently we only have one test, but if we expand our specification (we will), it doesn't make sense for the Driver to redial for every RPC call.

```go
package grpcserver

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Driver struct {
	Addr string

	connectionOnce sync.Once
	conn           *grpc.ClientConn
}

func (d *Driver) Greet(name string) (string, error) {
	conn, err := d.getConnection()
	if err != nil {
		return "", err
	}

	client := NewGreeterClient(conn)
	greeting, err := client.Greet(context.Background(), &GreetRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}

	return greeting.Message, nil
}

func (d *Driver) getConnection() (*grpc.ClientConn, error) {
	var err error
	d.connectionOnce.Do(func() {
		d.conn, err = grpc.Dial(d.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	})
	return d.conn, err
}
```

Here we're showing how we can use [`sync.Once`](https://pkg.go.dev/sync#Once) to ensure our `Driver` only attempts to create a connection to our server once.

Let's take a look at the current state of our project structure before moving on.

```
quii@Chriss-MacBook-Pro go-specs-greet % tree
.
├── adapters
│   ├── docker.go
│   ├── grpcserver
│   │   ├── driver.go
│   │   ├── greet.pb.go
│   │   ├── greet.proto
│   │   ├── greet_grpc.pb.go
│   │   └── server.go
│   └── httpserver
│       ├── driver.go
│       └── handler.go
├── cmd
│   ├── grpcserver
│   │   ├── Dockerfile
│   │   ├── greeter_server_test.go
│   │   └── main.go
│   └── httpserver
│       ├── Dockerfile
│       ├── greeter_server_test.go
│       └── main.go
├── go.mod
├── go.sum
├── greet.go
├── greet_test.go
└── specifications
    └── greet.go
```

- Adapters have cohesive units of functionality grouped together
- cmd holds our applications and acceptance tests in a very consistent structure
- Our domain code lives at the root, totally decoupled from any accidental complexity

### Consolidating `Dockerfile`

You've probably noticed the two `Dockerfiles` are almost identical beyond the path to the binary we wish to build.

`Dockerfiles` can accept arguments to let us re-use them in different contexts, which sounds perfect for us. We can delete our 2 Dockerfiles and instead have one at the root of the project with the following

```go
FROM golang:1.18-alpine

WORKDIR /app

ARG bin_to_build

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o svr cmd/${bin_to_build}/main.go

EXPOSE 50051
CMD [ "./svr" ]
```

We'll have to update our `StartDockerServer` function to pass in the argument when we build the images

```go
func StartDockerServer(
	ctx context.Context,
	t testing.TB,
	port string,
	binToBuild string,
) {
	t.Helper()
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../.",
			Dockerfile: "Dockerfile",
			BuildArgs: map[string]*string{
				"bin_to_build": &binToBuild,
			},
		},
		ExposedPorts: []string{fmt.Sprintf("%s:%s", port, port)},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(5 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})
}
```

And finally, update our tests to pass in the image to build (do this for the other test and change `grpcserver` to `httpserver`).

```go
func TestGreeterServer(t *testing.T) {
	var (
		ctx    = context.Background()
		port   = "50051"
		addr   = fmt.Sprintf("localhost:%s", port)
		driver = grpcserver.Driver{Addr: addr}
	)

	adapters.StartDockerServer(ctx, t, port, "grpcserver")
	specifications.GreetSpecification(t, &driver)
}
```

### Separating out different kinds of tests

Acceptance tests are great in that they test the whole system works from a pure user-facing, behavioural POV, but they do have their downsides compared to unit tests:

- Slower
- Quality of feedback is often not as focused as a unit test
- Doesn't help you with internal quality, or design

[The Test Pyramid](https://martinfowler.com/articles/practical-test-pyramid.html) gives us the guidance around the kind of mix we want for our test suite, you should read Fowler's post for more detail but the very simplistic summary for this post is "lots of unit tests, and a few acceptance tests".

For that reason, as a project grows you often may be in situations where the acceptance tests can take a few minutes to run. To offer a friendly developer experience for people checking out your project, you can enable developers to run the different kinds of tests separately.

It's preferable that running `go test ./...` should be runnable with no further set up from an engineer, beyond say a few key dependencies such as the Go compiler (obviously) and perhaps Docker.

Go provides a mechanism for engineers to run only "short" tests with the [short flag](https://pkg.go.dev/testing#Short)

`go test -short ./...`

We can add to our acceptance tests to see if the user wants to run our acceptance tests by inspecting the value of the flag

```go
if testing.Short() {
  t.Skip()
}
```

For this project I made a `Makefile` to show this usage

```makefile
build:
	golangci-lint run
	go test ./...

unit-tests:
	go test -short ./...

```

## Iterating on our work

With all this effort, you'd hope extending our system will now be simple. Making a system that is simple to work on, is not neccessairily easy, but it's worth the time, and is substantially easier to do when you start a project.

Let's extend our API to include a "curse" functionality.

## Write the test first

In our specification file, add the following

```go
type MeanGreeter interface {
	Curse(name string) (string, error)
}

func CurseSpecification(t *testing.T, meany MeanGreeter) {
	got, err := meany.Curse("Chris")
	assert.NoError(t, err)
	assert.Equal(t, got, "Go to hell, Chris!")
}
```

Pick one of our acceptance tests and try to use the specification

## Try to run the test



## Write the minimal amount of code for the test to run and check the failing test output

## Write enough code to make it pass

## Refactor

# Learn Go with Tests - Scaling Acceptance Tests (and light intro to gRPC)

This is a follow up to [Intro to acceptance tests](https://quii.gitbook.io/learn-go-with-tests/testing-fundamentals/intro-to-acceptance-tests)

## Ideas / things left to write about

- Discuss Dave Farleys acceptance test youtube vid, and reference gopherconuk talk by Riya and I
- Don't write an acceptance test for everything, reference the test pyramid again
  - Adding language switch would demand a change in the spec as we're changing the API of the system
  - Subsequent languages should be done via unit tests




## Things reader will learn

- How to use specifications / drivers to decouple the accidental and essential complexity
  - Normally when you're solving someone's problem, you're dealing with essential complexity, try and express that in the specification
- Top-down [GOOS](http://www.growing-object-oriented-software.com)-thinking
  - Start with hello, world, build from there
- Intro to gRPC



## Let's go

Create a new project

`go mod init github.com/quii/go-specs-greet` (replace `quii` with whatever you want, any imports in examples will be using this module though, adjust as necessary)

Make a folder `specifications` to hold our specification, and add a file `greet.go`

```go
package specifications

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

type Greeter interface {
	Greet() (string, error)
}

func GreetSpecification(t testing.TB, greeter Greeter) {
	got, err := greeter.Greet()
	assert.NoError(t, err)
	assert.Equal(t, got, "Hello, world")
}
```

My IDE (Goland) takes care of the fuss of adding dependencies for me, but if you need to do it manually you'd do

`go get github.com/alecthomas/assert/v2`

Given Farley's acceptance test design, we now have a specification which is decoupled from implementation. Interfaces are a great way to decouple code from implementation detail. The specification doesn't know, or care about _how_ we `Greet`, it's just concerned with the behaviour. This "behaviour" isn't much right now, but we'll expand upon the spec to add more functionality as we further iterate.

At this point, this level of ceremony to decouple our specification from implementation might make some people accuse us of "overly abstracting"; I promise you that acceptance tests that are too coupled to implementation become a real burden on engineering teams. I am confident to assert that most acceptance tests out in the wild are expensive to maintain, due to this inappropriate coupling; rather than the reverse, of being overly abstract.

We can use this specification to verify any "system" that can `Greet`.

### First system: HTTP API

Our requirement is to provider a greeter service over HTTP. So we'll need to create:

1. A **driver**. In this case, the way one works with a HTTP system is using a **HTTP client**. This code will know how to work with our API. Drivers implement the interface that specifications define.
2. A HTTP server with a greet API
3. A black-box, **acceptance test**, which is responsible for managing the life-cycle of spinning up the application, and then plugging the driver into the specification to run it as a test

## Write the test first

The initial process for creating an acceptance test that compiles and runs your program, executes the test and then cleans everything up can be quite labour intensive. Much less intensive in the long-run that repeatedly having to these steps yourself to check behaviour though!

It's preferable to set this up at the _start_ of your project on a very small amount of functionality. I typically start all my projects with a "hello world" server implementation, with all of my tests set up, ready for me to build the real functionality easily. Trying to retrofit acceptance testing into an existing system can be challenging, and without the setup existing for a project, engineers can tend to be lazy and not bother at all.

Most development teams these days are shipping using Docker, so our acceptance tests will test a docker image we'll build of our program.

To help us use Docker in our tests, we're going to use [Testcontainers](https://golang.testcontainers.org).

`go get github.com/testcontainers/testcontainers-go`

Create some structure to house our program we intend to ship

`mkdir -p cmd/http_server`

Inside the new folder, create a new file and add the following

`greeter_http_server_test.go`

```go
package main_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	go_specs_greet "github.com/quii/go-specs-greet"
	"github.com/quii/go-specs-greet/specifications"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGreeterServer(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../.",
			Dockerfile: "./cmd/http_server/Dockerfile",
		},
    ExposedPorts: []string{"8080:8080"},
		WaitingFor:   wait.ForHTTP("/").WithPort("8080"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	driver := go_specs_greet.Driver{BaseURL: "http://localhost:8080"}
	specifications.GreetSpecification(t, driver)
}
```

Notes:

- Most of the code is dedicated to building the Docker image of our web server and then launching a container from it
- We're allowing our driver to be configurable with the `BaseURL` field. This'll allow us to re-use the driver in other environments, such as staging, or even production.

## Try to run the test

```
./greeter_server_test.go:46:12: undefined: go_specs_greet.Driver
```

We're still practicing TDD here! It's a big first step we have to make, we need to make a few files and write maybe more code than we're typically used to, but when you're first starting this is often the case. It's so important we try and remember the rules of the red step.

> Commit as many sins as neccessary to get the test passing

## Write the minimal amount of code for the test to run and check the failing test output

Hold your nose, and remember we can refactor when the test is passing. Here's the code for our driver in `driver.go`

```go
package go_specs_greet

import (
	"io"
	"net/http"
)

type Driver struct {
	BaseURL string
}

func (d Driver) Greet() (string, error) {
	res, err := http.Get(d.BaseURL + "/greet")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	greeting, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(greeting), nil
}
```
The pattern to observe here is `Driver` implements the `Greeter` interface the specification defines. If you wish to re-use the specification against another system, all you have to do is follow this same pattern;
1. Implement the interface
2. Put the system specific code to drive the system for the test

Notes:

- You could argue that perhaps I should be writing tests to drive out the various `if err != nil`, but in my experience so long as you're not doing anything with the `err`, tests that say "you return the error you get" are fairly low value. With drivers you're typically going to be just bubbling up the error to the test anyway.
- **You shouldn't use the default HTTP client**. Later we'll pass in a HTTP client so it can be configured with timeouts e.t.c., but for now we're just trying to get ourselves to a passing test

Try and run the tests again, they should now compile, but not pass.

```
=== RUN   TestGreeterHandler
2022/09/10 18:49:44 Starting container id: 03e8588a1be4 image: docker.io/testcontainers/ryuk:0.3.3
2022/09/10 18:49:45 Waiting for container id 03e8588a1be4 image: docker.io/testcontainers/ryuk:0.3.3
2022/09/10 18:49:45 Container is ready id: 03e8588a1be4 image: docker.io/testcontainers/ryuk:0.3.3
    greeter_server_test.go:32: Did not expect an error but got:
        Error response from daemon: Cannot locate specified Dockerfile: ./cmd/http_server/Dockerfile: failed to create container
--- FAIL: TestGreeterHandler (0.59s)

```

We need to create a Dockerfile for our program. Inside our `http_server` folder, create a `Dockerfile` and add the following

```dockerfile
FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o svr cmd/http_server/*.go

EXPOSE 8080
CMD [ "./svr" ]
```

Don't worry too much about the details here, it can be refined and optimised, but for this example, it'll suffice. The advantage of our approach here is we can later improve our Dockerfile and have a test to prove it works as we intend it to. This is the real strength of having black-box tests!

Try and run the test again and it should complain about not being able to build the image. That's because we haven't added a program yet!

For the test to fully execute, we'll need to create a program that listens on `8080`, but **that's all**. Stick to the TDD discipline, don't write the production code that would make the test pass until we've verified the test fails as we'd expect.

Create a `main.go` inside our `http_server` folder with the following

```go
func main() {
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	})
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
```

```
    greet.go:16: Expected values to be equal:
        +Hello, World
        \ No newline at end of file
--- FAIL: TestGreeterHandler (2.09s)
```

## Write enough code to make it pass

Update the handler to behave how our specification wants it to

```go
func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "Hello, world")
	})
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
```

## Refactor

Whilst this technically isn't a refactor, we shouldn't rely on the default HTTP client, so let's change our client so it can be supplied one; which our test will give to it.

```go
type Driver struct {
	BaseURL string
	Client *http.Client
}

func (d Driver) Greet() (string, error) {
	res, err := d.Client.Get(d.BaseURL + "/greet")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	greeting, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(greeting), nil
}
```

Update the creation of the driver to pass in a client.

```go
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	driver := go_specs_greet.Driver{BaseURL: "http://localhost:8080", Client: &client}
	specifications.GreetSpecification(t, driver)
}
```

It's good practice to keep `main.go` as simple as possible, it only really aught to be concerned with piecing together the building blocks you make in to an application.

Create a file called `handler.go` and move our code into there

```go
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world")
}
```

Update `main.go` to import and use the handler instead

```go
package main

import (
	"net/http"

	go_specs_greet "github.com/quii/go-specs-greet"
)

func main() {
	handler := http.HandlerFunc(go_specs_greet.Handler)
	http.ListenAndServe(":8080", handler)
}
```

## Reflect

The first step felt like a fair amount of effort. We've made a number of `go` files to create and test a HTTP handler that returns a hard-coded string. This "iteration 0" ceremony and setup though will serve us well for further iterations.

Adding or changing functionality should be simple, and controlled by driving it through the specification and dealing with whatever changes it drives us to do. Now the `DockerFile` and `testcontainers` are set up for our acceptance test, we shouldn't have to change these files unless the way we construct our application changes.

We'll see this with our next requirement, greet a particular person.

## Write the test first

Edit our specification

```go
package specifications

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

type Greeter interface {
	Greet(name string) (string, error)
}

func GreetSpecification(t testing.TB, greeter Greeter) {
	got, err := greeter.Greet("Mike")
	assert.NoError(t, err)
	assert.Equal(t, got, "Hello, Mike")
}

```

To allow us to greet specific people, we need to change the interface to our system to accept a `name` parameter.

## Try to run the test

```
./greeter_server_test.go:48:39: cannot use driver (variable of type go_specs_greet.Driver) as type specifications.Greeter in argument to specifications.GreetSpecification:
	go_specs_greet.Driver does not implement specifications.Greeter (wrong type for Greet method)
		have Greet() (string, error)
		want Greet(name string) (string, error)
```

The change in the specification has meant our driver needs to be updated.

## Write the minimal amount of code for the test to run and check the failing test output

```go
func (d Driver) Greet(name string) (string, error) {
	res, err := d.Client.Get(d.BaseURL + "/greet?name=" + name)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	greeting, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(greeting), nil
}
```

The test should now run

```
    greet.go:16: Expected values to be equal:
        -Hello, world
        \ No newline at end of file
        +Hello, Mike
        \ No newline at end of file
--- FAIL: TestGreeterHandler (1.92s)
```

## Write enough code to make it pass

```go
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s", r.URL.Query().Get("name"))
}
```

## Refactor

In [HTTP Handlers Revisited](https://github.com/quii/learn-go-with-tests/blob/main/http-handlers-revisited.md) we discussed how important it is for HTTP handlers should only be response for handling HTTP concerns, any kind of "domain logic" should live outside of the handler. This allow us to develop domain logic in isolation of HTTP, making it simpler to test and understand.

Let's pull apart these concerns.

```go
func Handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprint(w, Greet(name))
}
```

In `greet.go`

```go
func Greet(name string) string {
	return fmt.Sprintf("Hello, %s", name)
}
```

## A slight diversion in to the "adapter" design pattern

Now that we've separated our domain logic of greeting people into a separate function, we are now free to write unit tests for our greet function; certainly a lot simpler than testing it through a specification, that goes through a driver, that hits a web server, to finally get a string!

Wouldn't it be nice if we could re-use our specification here too? After-all, the point of the specification is it's decoupled from implementation details.

Let's give it a go in `greet_test.go`

```go
func TestGreet(t *testing.T) {
	specifications.GreetSpecification(t, go_specs_greet.Greet)
}
```

This would be nice, but it doesn't work

```
./greet_test.go:11:39: cannot use go_specs_greet.Greet (value of type func(name string) string) as type specifications.Greeter in argument to specifications.GreetSpecification:
	func(name string) string does not implement specifications.Greeter (missing Greet method)
```

Our specification wants something that has a method `Greet()` not a function.

This is frustrating, we have a thing that we "know" is a `Greeter`, but it's not quite in the right **shape** for the compiler to let us use it. This is what the **adapter** pattern caters for.

> In [software engineering](https://en.wikipedia.org/wiki/Software_engineering), the **adapter pattern** is a [software design pattern](https://en.wikipedia.org/wiki/Software_design_pattern) (also known as [wrapper](https://en.wikipedia.org/wiki/Wrapper_function), an alternative naming shared with the [decorator pattern](https://en.wikipedia.org/wiki/Decorator_pattern)) that allows the [interface](https://en.wikipedia.org/wiki/Interface_(computer_science)) of an existing [class](https://en.wikipedia.org/wiki/Class_(computer_science)) to be used as another interface.[[1\]](https://en.wikipedia.org/wiki/Adapter_pattern#cite_note-HeadFirst-1) It is often used to make existing classes work with others without modifying their [source code](https://en.wikipedia.org/wiki/Source_code).

This is a lot of fancy words, for something that is quite simple. Which is often the case with design patterns, which is why people tend to roll their eyes at them. The value of design patterns is not specific implementations, but a language to describe certain solutions to common problems engineers face. If you have a team that has a shared vocabulary, it reduces the friction in communication.

Adapters allow you to "adapt" things to fit into other parts of your system.

Add this code in `greet.go`

```go
type GreetAdapter func(name string) string

func (g GreetAdapter) Greet(name string) (string, error) {
	return g(name), nil
}
```

We can now use our adapter in our test to plug our `Greet` function into the specification.

```go
func TestGreet(t *testing.T) {
	specifications.GreetSpecification(
		t,
		gospecsgreet.GreetAdapter(gospecsgreet.Greet),
	)
}
```

## Reflect

This felt simple right? OK, maybe it was simple due to the nature of the problem, but this method of work gives you discipline, a simple repeatable way of designing your code from top to bottom.

- Analyse your problem and identify a small improvement to your system that pushes you in the right direction
- Change the spec
- Follow the compilation errors until the test runs
- Update your implementation
- Refactor

After the pain of the first iteration, we didn't have to edit our acceptance test code at all because we have the seperation of specifications, drivers and implementation. Changing our specification required us to update our driver, and finally our implementation; but the boilerplate code around _how_ to spin up the system as a contaiiner was unaffected.

Even with the overhead of building a docker image for our application, and spinning up the container, the feedback loop for testing our **entire** application is very tight:

```
quii@Chriss-MacBook-Pro go-specs-greet % go test ./...
ok  	github.com/quii/go-specs-greet	0.181s
ok  	github.com/quii/go-specs-greet/cmd/httpserver	2.221s
?   	github.com/quii/go-specs-greet/specifications	[no test files]
```

Now, imagine your CTO has now decided gRPC is _the future_. She wants you to expose this same functionality over a gRPC server, whilst maintaining the existing HTTP server.

This is an example of **accidental complexity**. Accidental complexity is the complexity we have to deal with because we're working with computers, stuff like networks, disks, APIs, e.t.c. **Essential complexity** is sometimes referred to as "domain logic", it's the inescapable rules and truths within the domain you work in. They should be expressable to a non-technical person, and it's valuable to model them in our systems both in **specifications** and **domain code, that is decoupled from accidental complexity**. Many repository structures and design patterns are mainly dealing with this concern. For instance "ports and adapters" asks that you separate out your domain code from anything to do with accidental complexity, that code lives in an "adapters" folder.

Sometimes, it makes sense to do some refactoring before making a change

> Make the change easy, then make the change.

For that reason, let's gather our `http` code into a package called `httpserver` within an `adapters` folder

```
quii@Chriss-MacBook-Pro go-specs-greet % tree
.
├── adapters
│   └── httpserver
│       ├── driver.go
│       └── handler.go
├── cmd
│   └── httpserver
│       ├── Dockerfile
│       ├── greeter_server_test.go
│       └── main.go
├── go.mod
├── go.sum
├── greet.go
├── greet_test.go
└── specifications
    └── greet.go
```

Our domain code, our **essential complexity** lives at the root of our go module, and code that will allow us to use them in "the real world" are organised in to **adapters**. The `cmd` folder is where we can compose these logical groupings into useful applications, which have black-box tests to verify it all works. Nice!

Finally, we can do a _tiny_ bit of tidying up of our acceptance test. If you consider the high-level steps of our acceptance test:

- Build _some_ docker image
- Wait for it to be listening on _some_ port
- Create _some_ driver to send messages to that port
- Plug in the driver into the specification

... you'll realise we have the same requirements for an acceptance test for the gRPC server!

The `adapters` folder seems a good a place as any, so inside a file called `docker.go` , encapsulate the first 2 steps in a function that we'll re-use next.

```go
package adapters

import (
	"context"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartDockerServer(
  ctx context.Context,
	t testing.TB,
	dockerFilePath string,
	port string,
) {
	t.Helper()
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../.",
			Dockerfile: dockerFilePath,
		},
		ExposedPorts: []string{fmt.Sprintf("%s:%s", port, port)},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(5 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})
}
```

This gives us an opportunity to clean up our acceptance test a little

```go
func TestGreeterServer(t *testing.T) {
	var (
		ctx            = context.Background()
		port           = "8080"
		dockerFilePath = "./cmd/httpserver/Dockerfile"
		baseURL        = fmt.Sprintf("http://localhost:%s", port)
		driver         = go_specs_greet.Driver{BaseURL: baseURL, Client: &http.Client{
			Timeout: 1 * time.Second,
		}}
	)

	adapters.StartDockerServer(ctx, t, dockerFilePath, port)
	specifications.GreetSpecification(t, driver)
}
```

This should make writing the _next_ test simpler.

## Write the test first

You can imagine this functionality being a new adapter in to our domain code. For that reason we:

- Shouldn't have to change the specification;
- Should be able to re-use the specification;
- Should be able to re-use the domain code.

Create a new folder `grpcserver` inside `cmd` to house our new program and the corresponding acceptance test. Inside `cmd/grpc_server/greeter_server_test.go` add an acceptance test, which, not by coincedence, but by design, looks very similar to our HTTP server test.

```go
package main_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/quii/go-specs-greet/adapters"
	"github.com/quii/go-specs-greet/adapters/grpcserver"
	"github.com/quii/go-specs-greet/specifications"
)

func TestGreeterServer(t *testing.T) {
	var (
		ctx            = context.Background()
		port           = "50051"
		dockerFilePath = "./cmd/grpcserver/Dockerfile"
		addr           = fmt.Sprintf("localhost:%s", port)
		driver         = grpcserver.Driver{Addr: addr}
	)

	adapters.StartDockerServer(ctx, t, dockerFilePath, port)
	specifications.GreetSpecification(t, &driver)
}
```

The only differences are:

- We use a different docker file, because we're building a different program
- We use a different driver to plug in to the specification

## Try to run the test

```
./greeter_server_test.go:26:12: undefined: grpcserver
```

We haven't created a Driver yet, so it won't compile.

## Write the minimal amount of code for the test to run and check the failing test output

Create a `grpcserver` folder inside `adapters` and inside it create `driver.go`

```go
package grpcserver

type Driver struct {
	Addr string
}

func (d Driver) Greet(name string) (string, error) {
	return "", nil
}
```

If you run again, it should now _compile_ but not pass, because we haven't created a Dockerfile and corresponding program for it to run against.

Create a new `Dockerfile` inside `cmd/grpcserver`.

```go
FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o svr cmd/grpcserver/*.go

EXPOSE 8080
CMD [ "./svr" ]
```

And a `main.go`

```go
package main

import "fmt"

func main() {
	fmt.Println("implement me")
}
```

You should find now that the test fails because our server is not listening on the port. Now is the time to start building our client and server with gRPC

## Write enough code to make it pass

### gRPC

If you're unfamiliar with gRPC, I'd start by looking at the [gRPC website](https://grpc.io) but for the purposes of this chapter, it's just another kind of adapter in to our system, a way of other systems being able to call (**r**emote **p**rocedure **c**all) our amazing domain code.

The twist is you define a "service definition" using Protocol Buffers. You can then generate server and client code from the definition. This not only works for Go, but for most mainstream languages too. This means you can share a definition with other teams in your company who may not even write Go, and still be able to do service to service communucation very smoothly.

If you haven't used gRPC before you'll need to install a **Protocol buffer compiler** and some **Go plugins** for it. [The gRPC website has clear instructions as to how to do this](https://grpc.io/docs/languages/go/quickstart/).

Inside the same folder as our new driver, add a `greet.proto` file with the following

```protobuf
syntax = "proto3";

option go_package = "github.com/quii/adapters/grpcserver";

package grpcserver;

service Greeter {
  rpc Greet (GreetRequest) returns (GreetReply) {}
}

message GreetRequest {
  string name = 1;
}

message GreetReply {
  string message = 1;
}
```

You don't need to be an expert in Protocol Buffers to follow this definition. We're defining a service, which has a `Greet` method, and then describing the incoming and outgoing message types.

Inside `adapters/grpcserver` run the following to generate the client and server code

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    greet.proto
```

If it worked, we will have some code generated for us to use. Let's start by using the generated client code inside our `Driver`.

```go
package grpcserver

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Driver struct {
	Addr string
}

func (d Driver) Greet(name string) (string, error) {
	//todo: we shouldn't redial every time we call greet, refactor out when we're green
	conn, err := grpc.Dial(d.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := NewGreeterClient(conn)
	greeting, err := client.Greet(context.Background(), &GreetRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}

	return greeting.Message, nil
}

```

Now that we have a client, we need to update our `main.go` to create a server. Remember at this point we're just trying to get our test to pass, and not worrying about code quality.

```go
package main

import (
	"context"
	"log"
	"net"

	"github.com/quii/go-specs-greet/adapters/grpcserver"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	grpcserver.RegisterGreeterServer(s, &GreetServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

type GreetServer struct {
	grpcserver.UnimplementedGreeterServer
}

func (g GreetServer) Greet(ctx context.Context, request *grpcserver.GreetRequest) (*grpcserver.GreetReply, error) {
	return &grpcserver.GreetReply{Message: "fixme"}, nil
}
```

To create our gRPC server, we have to implement the interface it generated for us

```go
// GreeterServer is the server API for Greeter service.
// All implementations must embed UnimplementedGreeterServer
// for forward compatibility
type GreeterServer interface {
	Greet(context.Context, *GreetRequest) (*GreetReply, error)
	mustEmbedUnimplementedGreeterServer()
}
```

- Listen on the port
- We create a `GreetServer` that implements this interface, and then register it with `grpcServer.RegisterGreeterServer`, along with a `grpc.Server`.
- Use the server with the listener



It wouldn't be a huge extra effort to call our domain code inside `greetServer.Greet` rather than hard-coding `fix-me` in the message, but I'd like to run our acceptance test first, just to see everything is working end to end on a transport level

```
greet.go:16: Expected values to be equal:
-fixme
\ No newline at end of file
+Hello, Mike
\ No newline at end of file
```

Nice! We can see our driver is able to connect to our gRPC server in the test.

Finally, we can call our domain code inside our `GreetServer`

```go
type GreetServer struct {
	grpcserver.UnimplementedGreeterServer
}

func (g GreetServer) Greet(ctx context.Context, request *grpcserver.GreetRequest) (*grpcserver.GreetReply, error) {
	return &grpcserver.GreetReply{Message: gospecsgreet.Greet(request.Name)}, nil
}
```

Finally it passes! We have an acceptance test that proves our gRPC greet server behaves how we'd like.

## Refactor

We committed a number of sins to get the test passing, but now they're passing we have the safety-net to refactor.

### Simplify main

Like before, we don't want `main` having too much code inside it, and it feels inconsistent with our other implementation. We can move our new `GreetServer` into `adapters/grpcserver` as that's definitely where it should live. In terms of cohesion if we happen to change the service definition, we want the "blast-radius" of change to be confined to that area of our code.

### Don't redial in our driver every time

Currently we only have one test, but if we expand our specification (we will), it doesn't make sense for the Driver to redial for every RPC call.

```go
package grpcserver

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Driver struct {
	Addr string

	connectionOnce sync.Once
	conn           *grpc.ClientConn
}

func (d *Driver) Greet(name string) (string, error) {
	conn, err := d.getConnection()
	if err != nil {
		return "", err
	}

	client := NewGreeterClient(conn)
	greeting, err := client.Greet(context.Background(), &GreetRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}

	return greeting.Message, nil
}

func (d *Driver) getConnection() (*grpc.ClientConn, error) {
	var err error
	d.connectionOnce.Do(func() {
		d.conn, err = grpc.Dial(d.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	})
	return d.conn, err
}
```

Here we're showing how we can use [`sync.Once`](https://pkg.go.dev/sync#Once) to ensure our `Driver` only attempts to create a connection to our server once.

Let's take a look at the current state of our project structure before moving on.

```
quii@Chriss-MacBook-Pro go-specs-greet % tree
.
├── adapters
│   ├── docker.go
│   ├── grpcserver
│   │   ├── driver.go
│   │   ├── greet.pb.go
│   │   ├── greet.proto
│   │   ├── greet_grpc.pb.go
│   │   └── server.go
│   └── httpserver
│       ├── driver.go
│       └── handler.go
├── cmd
│   ├── grpcserver
│   │   ├── Dockerfile
│   │   ├── greeter_server_test.go
│   │   └── main.go
│   └── httpserver
│       ├── Dockerfile
│       ├── greeter_server_test.go
│       └── main.go
├── go.mod
├── go.sum
├── greet.go
├── greet_test.go
└── specifications
    └── greet.go
```

- Adapters have cohesive units of functionality grouped together
- cmd holds our applications and acceptance tests in a very consistent structure
- Our domain code lives at the root, totally decoupled from any accidental complexity

### Consolidating `Dockerfile`

You've probably noticed the two `Dockerfiles` are almost identical beyond the path to the binary we wish to build.

`Dockerfiles` can accept arguments to let us re-use them in different contexts, which sounds perfect for us. We can delete our 2 Dockerfiles and instead have one at the root of the project with the following

```go
FROM golang:1.18-alpine

WORKDIR /app

ARG bin_to_build

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o svr cmd/${bin_to_build}/main.go

EXPOSE 50051
CMD [ "./svr" ]
```

We'll have to update our `StartDockerServer` function to pass in the argument when we build the images

```go
func StartDockerServer(
	ctx context.Context,
	t testing.TB,
	port string,
	binToBuild string,
) {
	t.Helper()
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../.",
			Dockerfile: "Dockerfile",
			BuildArgs: map[string]*string{
				"bin_to_build": &binToBuild,
			},
		},
		ExposedPorts: []string{fmt.Sprintf("%s:%s", port, port)},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(5 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})
}
```

And finally, update our tests to pass in the image to build (do this for the other test and change `grpcserver` to `httpserver`).

```go
func TestGreeterServer(t *testing.T) {
	var (
		ctx    = context.Background()
		port   = "50051"
		addr   = fmt.Sprintf("localhost:%s", port)
		driver = grpcserver.Driver{Addr: addr}
	)

	adapters.StartDockerServer(ctx, t, port, "grpcserver")
	specifications.GreetSpecification(t, &driver)
}
```

### Separating out different kinds of tests

Acceptance tests are great in that they test the whole system works from a pure user-facing, behavioural POV, but they do have their downsides compared to unit tests:

- Slower
- Quality of feedback is often not as focused as a unit test
- Doesn't help you with internal quality, or design

[The Test Pyramid](https://martinfowler.com/articles/practical-test-pyramid.html) gives us the guidance around the kind of mix we want for our test suite, you should read Fowler's post for more detail but the very simplistic summary for this post is "lots of unit tests, and a few acceptance tests".

For that reason, as a project grows you often may be in situations where the acceptance tests can take a few minutes to run. To offer a friendly developer experience for people checking out your project, you can enable developers to run the different kinds of tests separately.

It's preferable that running `go test ./...` should be runnable with no further set up from an engineer, beyond say a few key dependencies such as the Go compiler (obviously) and perhaps Docker.

Go provides a mechanism for engineers to run only "short" tests with the [short flag](https://pkg.go.dev/testing#Short)

`go test -short ./...`

We can add to our acceptance tests to see if the user wants to run our acceptance tests by inspecting the value of the flag

```go
if testing.Short() {
  t.Skip()
}
```

For this project I made a `Makefile` to show this usage

```makefile
build:
	golangci-lint run
	go test ./...

unit-tests:
	go test -short ./...

```

### Primitive obsession

We want our specifications to be an abstract description of our domain, but we have coupled it directly to some very specific types. Namely, `strings`.

TODO: Some resources on primitive obsession.

## Iterating on our work

With all this effort, you'd hope extending our system will now be simple. Making a system that is simple to work on, is not neccessairily easy, but it's worth the time, and is substantially easier to do when you start a project.

Let's extend our API to include a "curse" functionality.

## Write the test first

In our specification file, add the following

```go
type MeanGreeter interface {
	Curse(name string) (string, error)
}

func CurseSpecification(t *testing.T, meany MeanGreeter) {
	got, err := meany.Curse("Chris")
	assert.NoError(t, err)
	assert.Equal(t, got, "Go to hell, Chris!")
}
```

Pick one of our acceptance tests and try to use the specification

```go
func TestGreeterServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var (
		ctx    = context.Background()
		port   = "50051"
		addr   = fmt.Sprintf("localhost:%s", port)
		driver = grpcserver.Driver{Addr: addr}
	)

	t.Cleanup(driver.Close)
	adapters.StartDockerServer(ctx, t, port, "grpcserver")
	specifications.GreetSpecification(t, &driver)
	specifications.CurseSpecification(t, &driver)
}
```

## Try to run the test

```
# github.com/quii/go-specs-greet/cmd/grpcserver_test [github.com/quii/go-specs-greet/cmd/grpcserver.test]
./greeter_server_test.go:27:39: cannot use &driver (value of type *grpcserver.Driver) as type specifications.MeanGreeter in argument to specifications.CurseSpecification:
	*grpcserver.Driver does not implement specifications.MeanGreeter (missing Curse method)

```

Our `Driver` doesn't support `Curse` yet.

## Write the minimal amount of code for the test to run and check the failing test output

Remember we're just trying to get the test to run, so add the method to `Driver`

```go
func (d *Driver) Curse(name string) (string, error) {
	return "", nil
}
```

If you try again, the test should compile, run, and fail

```
greet.go:26: Expected values to be equal:
+Go to hell, Chris!
\ No newline at end of file
```

## Write enough code to make it pass

We'll need to update our protocol buffer specification have a `Curse` method on it, and then regenerate our code.

```protobuf
service Greeter {
  rpc Greet (GreetRequest) returns (GreetReply) {}
  rpc Curse (GreetRequest) returns (GreetReply) {}
}
```

You could argue that re-using the types `GreetRequest` and `GreetReply` is inappropriate coupling, but we can deal with that in the refactoring stage. As I keep stressing, we're just trying to get the test passing so we verify the software works, _then_ we can make it nice.

Re-generate our code with (inside `adapters/grpcserver`).

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    greet.proto
```

### Update driver

Now the client code has been updated, we can now call `Curse` in our `Driver`

```go
func (d *Driver) Curse(name string) (string, error) {
	conn, err := d.getConnection()
	if err != nil {
		return "", err
	}

	client := NewGreeterClient(conn)
	greeting, err := client.Curse(context.Background(), &GreetRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}

	return greeting.Message, nil
}
```

### Update server

Finally, we need to add the `Curse` method to our `Server`

```go
package grpcserver

import (
	"context"
	"fmt"

	gospecsgreet "github.com/quii/go-specs-greet"
)

type GreetServer struct {
	UnimplementedGreeterServer
}

func (g GreetServer) Curse(ctx context.Context, request *GreetRequest) (*GreetReply, error) {
	return &GreetReply{Message: fmt.Sprintf("Go to hell, %s!", request.Name)}, nil
}

func (g GreetServer) Greet(ctx context.Context, request *GreetRequest) (*GreetReply, error) {
	return &GreetReply{Message: gospecsgreet.Greet(request.Name)}, nil
}
```

The tests should now pass.

## Refactor

Try doing this yourself.

- Extract the "domain logic", away from the grpc server, like we did for `Greet`. Use the specification as a unit test against your domain logic
- Have separate types in the protobuf to ensure the message types for `Greet` and `Curse` are decoupled.

## Implementing `Curse` for the HTTP server

Again, an exercise for you, the reader. We have our domain-level specification, and we have our domain-level logic neatly separated. If you've followed this chapter, this should be very straightforward.

- Add the specification to the existing acceptance test for the HTTP server
- Update your `Driver`
- Add the new endpoint to the server, and re-use the domain code to implement the functionality. You may wish to use `http.NewServeMux` to handle the routiing to the separate endpoints.

Remember to work in small steps, commit and run your tests frequently. If you get really stuck [you can find my implementation on GitHub](https://github.com/quii/go-specs-greet).

## Wrapping up

From here, hopefully you can see the predictable, structuted workflow for driving change on our application with this approach.

On your day job you can imagine talking to a stakeholder who wants to extend the system you work on in some way. Simply capture it in a domain-centric, implementation-agnostic way in the specification, and use it as a north-star towards your efforts. By separating out the concerns of essential complexity and accidental complexity, your work will feel less ad-hoc, and more structured and deliberate. 

