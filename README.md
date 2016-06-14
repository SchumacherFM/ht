HT: HTTP Testing Made Easy
==========================

[![Go Report Card](https://goreportcard.com/badge/github.com/vdobler/ht)](https://goreportcard.com/report/github.com/vdobler/ht)
[![GoDoc](https://godoc.org/github.com/vdobler/ht?status.svg)](https://godoc.org/github.com/vdobler/ht)

End-to-end testing of HTTP requests/responses is easy with Go.

Writing and maintaining high level test is even easier with ht.


Goals
-----

ht tries to achieve the following:
* Make generating all common HTTP requests easy.
* Provide high- and low-level checks on the received response.
* Make measuring response times easy.


Non-Goals
---------

ht is not the jack of all trades in testing web applications:
* Simulating different browsers is not done.
* Simulating interactions in a browser is not done.


Installation
------------

Installing ht should be simple if Go 1.5 and git are available and working:
* Run `GO15VENDOREXPERIMENT=1 go get github.com/vdobler/ht/cmd/ht`
  which should download, compile and install everything.
* Run `$GOPATH/bin/ht help` to get you started.
* For a quick check of a HTML page do a 

    `$GOPATH/bin/ht quick <URL-of-HTML-page>`

  and check the generated `Report.html` file.

If you want to run checks on rendered HTML pages you need a local installation
of PhantomJS in version >= 2.0. See http://phantomjs.org .


Documentation
-------------

For a start have a look at the 

Tutorial https://github.com/vdobler/ht/blob/master/cmd/ht/Tutorial.md

or the 

Showcase here https://github.com/vdobler/ht/tree/master/showcase

The showcase is pretty nonsensical but show almost all features
in just a few files.  You might want to `$ go run showcase.go` to
have a dummy server listening on localhost:8080 to run the tests
against: `$ ht exec showcase.suite`


For details see the the godoc for reference:

* Tests, Checks and their options:
  https://godoc.org/github.com/vdobler/ht/ht
* Details to Conditions used in Checks:
  https://godoc.org/github.com/vdobler/ht/condition
* Image fingerprinting:
  https://godoc.org/github.com/vdobler/ht/fingerprint
* The ht command itself:
  https://godoc.org/github.com/vdobler/ht/cmd/ht
  Run `ht help` for details
* An example test suite can be found in
  https://github.com/vdobler/ht/blob/master/testdata/sample.suite
