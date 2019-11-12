## 0x Mesh Browser Example

This directory contains an example of how to run Mesh in the browser.

### Running the Example

To run the example, first change into the __examples/browser__ directory and
then run:

```
yarn install --force && yarn build
```

Then simply serve the __examples/browser/dist__ directory using the web server
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
- [Browser API Documentation](https://0x-org.gitbook.io/mesh/getting-started/reference)
