[![Go Reference](https://pkg.go.dev/badge/github.com/kentzo/dsomessage.svg)](https://pkg.go.dev/github.com/kentzo/dsomessage)
[![Coverage Status](https://coveralls.io/repos/github/Kentzo/dsomessage/badge.svg?branch=main)](https://coveralls.io/github/Kentzo/dsomessage?branch=main)
[![CI](https://github.com/kentzo/dsomessage/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/kentzo/dsomessage/actions/workflows/ci.yml?query=branch%3Amain)

# dsomessage

dsomessage provides [RFC 8490][rfc8490] DNS Stateful Operations and [RFC 8765][rfc8765] DNS Push Notifications primitives
as well as parsing and building machinery.

The goal is minimal overhead and control over heap allocations and processing.

[rfc8490]: https://www.rfc-editor.org/rfc/rfc8490.html
[rfc8765]: https://www.rfc-editor.org/rfc/rfc8765.html
