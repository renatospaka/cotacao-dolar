# Go
This is a Go in Docker boilerplate I use myself whenever I start a Go project from scratch.

Enjoy it. Make suggestions if you like so.

# necessary steps
## check doors
- docker-compose.yaml ->
  - check if doors mapping is linking to a free port
  (8000:8000)

## 1st execution:
- comment lines 21 to 32 in Dockerfile
- remove go.mod e go.sum (if there are any)

## next executions:
- uncomment lines 21 to 32 in Dockerfile
- execute once -> docker-compose up --build
(to rebuilt the container with uncommented lines)
- start programing in Go

# go have fun using Go -> https://go.dev/solutions/#case-studies
