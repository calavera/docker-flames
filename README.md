# docker-flames

docker-flames is a suite of tests to run performance tests on a remote Docker daemon.

# How to

docker-flames uses Digital Ocean to provision several virtual machines and bootstrap the Docker Engine on them.
Then, it uses a set of Go programs to exercise the remote process and gather pprof information.

# License

MIT. See [LICENSE](LICENSE)
