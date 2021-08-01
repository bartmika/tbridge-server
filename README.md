## Usage 1 of 2: Server Mode

Make sure the [`tstorage-server`](https://github.com/bartmika/tstorage-server) is running. Please note the default address is `50051`.

```
$GOBIN/tstorage-server serve
```

Generate our environment variables for the server side by running the following.

```
go run main.go create_uuid
go run main.go create_hmacsecret
```

Please following the instructions that are printed in the terminal when running each sub-commands. Here is an example of something that you would run on the server side.

```
export TBRIDGE_SERVER_SESSION_UUID=50d2d42e-52f2-4a2b-87c2-62e96cb8522c
export TBRIDGE_SERVER_HMAC_SECRET=sjSCNSbfBGrVXwOvmAdasrFaqIYSRZFlrqFjJfUcpuNCQdezbFv
```

Next is you need to create an `authentication token` by running the following sub-command and save the output for later use.

```
go run main.go create_token
```

Start our bridge in `server mode` so it will connect to our `tstorage-server` and accept HTTP requests from the network.

```
go run main.go server_mode --storage_addr="localhost:50051"
```

## Usage 2 of 2: Client Mode

Paste the saved result when you ran the `create_token` sub-command. Here is an example of what you should see your console.

```
Run in your console:

export TBRIDGE_CLIENT_ACCESS_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjgzOTQ1NzAsInNlc3Npb25fdXVpZCI6IjUwZDJkNDJlLTUyZjItNGEyYi04N2MyLTYyZTk2Y2I4NTIyYyJ9._E66kVVy9c2gKU79fdrCY4IJck4Dpb6skrk4BmimdBw

export TBRIDGE_CLIENT_REFRESH_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mjg2NTM3NzAsInNlc3Npb25fdXVpZCI6IjUwZDJkNDJlLTUyZjItNGEyYi04N2MyLTYyZTk2Y2I4NTIyYyJ9.ki0OL7SnI45SKvCXbUyp99SPq42gOrxXCfflCGccnU4
```

Start our bridge in `client mode` so it will provide a gRPC interface for all local applications to use. All gRPC requests
made to this bridge will be converted into HTTP requests and sent over the network to the bridge running in `server mode`.

```
go run main.go client_mode --port=50053 --remote_addr="http://localhost:5000"
```

Our bridges are connected! Your local applications can make gRPC requests to the remote `tstorage-server` as if the storage server
was loaded locally.
