# File Watch example.

**NOTE:** This is a fork/vendoring of http://github.com/gorilla/websocket
The following documentation has been modified to point at this fork for
convenience.

This example sends a file to the browser client for display whenever the file is modified.

    $ go get github.com/btcsuite/websocket
    $ cd `go list -f '{{.Dir}}' github.com/btcsuite/websocket/examples/filewatch`
    $ go run main.go <name of file to watch>
    # Open http://localhost:8080/ .
    # Modify the file to see it update in the browser.
