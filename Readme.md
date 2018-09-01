# blackbox - Just a logger

Blackbox is a dead simple logger that will do exactly what you want a logger to
do; log messages; messages that are clear and contextual.

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"

    "github.com/RobertWHurst/blackbox"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
    logger := blackbox.New()
    logger.Info("Starting HTTP server")

    logger.Debug("Binding root handler")
    http.HandleFunc("/", handler)

    logger.Debug("Binding port")
    go func() {
      err := http.ListenAndServe(":8080", nil)
      logger.Fatal(err)
    }()
    logger.Info("Ready and awaiting connections")

    sigChan := make(chan os.Signal, 1)
    signal.Notify(
        sigChan,
        syscall.SIGHUP,
        syscall.SIGINT,
        syscall.SIGTERM,
        syscall.SIGQUIT,
    )
    
    <- sigChan
    logger.Debug("The operating system has requested process termination")
    logger.Info("Shutting down. Goodbye!")
    os.Exit(0)
}
```
