# dawa
Golang implementation of DAWA AWS Suite 4 (Danish Address Info)

The dawa package can be used to de-serialize structures received from "Danmarks Adressers Web API (DAWA)" (Addresses of Denmark Web API).

* This package allows to de-serialize JSON responses from the web api into typed structs.
* The package also allows importing JSON or CSV downloads from the official web page.

See the /examples folder for more information.

Package home: https://github.com/klauspost/dawa

Information abou the format and download/API options, see http://dawa.aws.dk/

# Installation 

```go get github.com/klauspost/dawa/...```

This will also install the only dependecy "go-codec": https://github.com/ugorji/go which is used to decode JSON streams more efficiently than the standard golang libraries.

To use the library in your own code, simply add ```"github.com/klauspost/dawa"``` to your imports.

# Documentation
[![GoDoc][1]][2] [![Build Status][3]][4]
[1]: https://godoc.org/github.com/klauspost/dawa?status.svg
[2]: https://godoc.org/github.com/klauspost/dawa
[3]: https://travis-ci.org/klauspost/dawa.svg
[4]: https://travis-ci.org/klauspost/dawa

# Description text in Danish:

Danmarks Adressers Web API (DAWA) udstiller data og funktionalitet vedrørende Danmarks adresser, adgangsadresser, vejnavne samt postnumre.
DAWA anvendes til etablering af adressefunktionalitet i it-systemer. Målgruppen for nærværende website er udviklere, som ønsker at indbygge adressefunktionalitet i deres it-systemer.
