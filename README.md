# testflight demo

This is a demo of [testflight](https://github.com/drewolson/testflight) for the [Minneapolis Go Meetup](http://www.meetup.com/golangmn/).

Here's the environment variables I set before doing anything.
```bash
export GOPATH=/home/alex/src/testflight-demo
export PATH="$GOPATH/bin:$PATH"
export TEST_DATADIR=/home/alex/src/testflight-demo/src/github.com/waucka/testflight-demo/testdata
```

You will also need a MongoDB server.  Tell the program where the server
is by setting `MONGO_HOST`.

Additionally, you will need `gocov`.  I recommend using `go get` to install it.