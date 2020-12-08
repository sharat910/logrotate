# logrotator

A time-based rotating log writer for Golang. 
`LogRotator` implements the standard `io.WriteCloser` interface and can be plugged into other logging packages or can be used to write rotating csvs, jsons or other files. 

Let's say that a golang application writes continuous logs into `logs/example.log`. 
It can use the `logrotate` package to split up the logs every (configurable) `rotateInterval` (defaults to 1 day) to create 
multiple log files (e.g. `logs/2020-12-08-example.log`, `logs/2020-12-09-example.log`, and so on.).

Example usage: `cmd/example/main.go`

#### TODO
* Add more tests
* Compress old logs