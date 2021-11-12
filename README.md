# Portal

Portal is a domain name server written in Go. Portal is used by [Void](https://github.com/go-void/void) to answer (and
block) DNS queries.

## General TODOs

-   **Logging:** At which level should we log and how do we pass the logger around in the best / most efficient way?
-   **Data collection:** How should we collect metrics / statistics which can be shown in the web interface of void?

## RFC TODOs

-   IN-ADDR.ARPA domain [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035#section-3.5)
-   Message compression [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.4)
-   Master Files (Lexing + Parsing) [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035#section-5)
-   Implement RPZs [DRAFT Vixie DNS RPZ](https://datatracker.ietf.org/doc/html/draft-vixie-dns-rpz-00)

## RFS to implement

-   LOC RR [RFC 1876](https://datatracker.ietf.org/doc/html/rfc1876)
-   Implement EDNS [RFC 6891](https://datatracker.ietf.org/doc/html/rfc6891)

## Supported RFCs

- Domain Names - Concepts and Facilities [RFC 1034](https://datatracker.ietf.org/doc/html/rfc1034)
- Domain Names - Implementation and Specification [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035)
- Serial Number Arithmetic [RFC 1982](https://datatracker.ietf.org/doc/html/rfc1982)