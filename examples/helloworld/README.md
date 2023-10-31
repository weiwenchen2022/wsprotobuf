# Hello Example

This directory contains a greet server example using github.com/weiwenchen2022/wsprotobuf.

1. Run the server:

```bash
$ cd examples/helloworld
$ go run ./server
listening on http://127.0.0.1:51055
```

3. Run the client:

```bash
$ cd examples/helloworld
$ go run ./client ws://127.0.0.1:51055
Greeting: Hello world
```
