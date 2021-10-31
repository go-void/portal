# Portal

Portal is a domain name server written in Go. Portal is used by [Void](https://github.com/go-void/void) to answer (and
block) DNS queries.

## TODOs

-   IN-ADDR.ARPA domain [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035#section-3.5)
-   Message compression [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.4)
-   Master Files (Lexing + Parsing) [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035#section-5)
-   Implement RPZs [DRAFT Vixie DNS RPZ](https://datatracker.ietf.org/doc/html/draft-vixie-dns-rpz-00)
-   Implement EDNS [RFC 6891](https://datatracker.ietf.org/doc/html/rfc6891)

-   Rework cache / store / tree
-   Logging
