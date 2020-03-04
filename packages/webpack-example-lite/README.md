## 0x Mesh Browser Lite Example

This directory contains an example of how to run Mesh in the browser by serving
raw WebAssembly bytecode directly. This manner of serving Mesh to the browser
can have performance benefits and allow the bundle size of an application to remain
much smaller.

### Running the Example

To run the example, first build the monorepo by changing into the __../../__
directory (the project's root directory) and then run:

```
yarn install && yarn build
```

Then simply serve the __./dist__ directory using the web server
of your choice and navigate to the page in your browser. For example, you could
use `goexec`:

```
go get -u github.com/shurcooL/goexec
goexec 'http.ListenAndServe(":8000", http.FileServer(http.Dir("./dist")))'
```

Then navigate to [localhost:8000](http://localhost:8000) in your browser (Chrome
or Firefox are recommended).

### More Information

- [Browser Guide](https://0x-org.gitbook.io/mesh/getting-started/browser)
- [Browser Lite API Documentation](https://0x-org.gitbook.io/mesh/getting-started/browser-lite/reference)
