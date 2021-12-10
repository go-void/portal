# Portal

Portal is a domain name server written in Go. Portal is used by [Void](https://github.com/go-void/void) to answer (and
block) DNS queries.

## TODOs

### General

- At which level should we log and how do we pass the logger around in the best / most efficient way?
- How should we collect metrics / statistics which can be shown in the web interface of void?
- Implement mechanism to return FORMERR when unpacking a message (prior to handling). Maybe move the unpacking into the
  handling func to be able to handle these kind of errors better and in a centralized place. (See OPT record)
- Figure out a way to handle EDNS Options with access to some core components like cache, etc.

#### Caching 

We need to support lookup of partial domain names: We currently do the following:

- Check if we cached the whole domain name, e.g. `example.com`
- If we did, return requested RR

If we recursively resolve a domain name and encounter a NS record without any glue records we always start to resolve
the domain name of the NS record from root (.). Instead we should lookup the cache if we have partial stored records of
NS RRs.

### RFCs

-   IN-ADDR.ARPA domain [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035#section-3.5)
-   Message compression [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.4)
-   Finalize Master Files Parsing [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035#section-5)
-   Implement RPZs [DRAFT Vixie DNS RPZ](https://datatracker.ietf.org/doc/html/draft-vixie-dns-rpz-00)

### RFS to implement

- LOC RR [RFC 1876](https://datatracker.ietf.org/doc/html/rfc1876)
- Incremental Zone Transfer in DNS [RFC 1995](https://datatracker.ietf.org/doc/html/rfc1995)
- Implement EDNS [RFC 6891](https://datatracker.ietf.org/doc/html/rfc6891)
  - Apple's DNS Long-Lived Queries Protocol [RFC 8764](https://datatracker.ietf.org/doc/html/rfc8764)
  - DNS Name Server Identifier (NSID) [RFC 5001](https://datatracker.ietf.org/doc/html/rfc5001)

## Supported RFCs

- Domain Names - Concepts and Facilities [RFC 1034](https://datatracker.ietf.org/doc/html/rfc1034)
- Domain Names - Implementation and Specification [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035)
- Serial Number Arithmetic [RFC 1982](https://datatracker.ietf.org/doc/html/rfc1982)

## Usage

### Standalone

To use the server as a standalone DNS server, follow these steps:

1. Clone the repository
2. Build the binary via `go build`
3. Run the binary with `./portal`

### Library

```go
import (
    "fmt"
    "os"

    "github.com/go-void/portal/pkg/config"
    "github.com/go-void/portal/pkg/server"
)

func main() {
    // Create new server
    s := server.New(&config.Config{})

    // Optionally provide custom implementations for different components via
    err := s.Configure()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // Run the server
    err = s.Run()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // This blocks until the server is shutdown
    s.Wait()
}
```