go_http_multiplexer
===================

A web service written in Go that fetches multiple URLs concurrently and returns them encoded in JSON

Running
-------
Either

    ./server

or

    go run server.go

and it will be running on port 12345.

Using
-----

Send a `GET` request to `/multiget` with the `url` parameter specified for each URL you
want to fetch.

Example:

    curl -v "http://localhost:12345/multiget?url=http://jak.io&url=http://samscupcakery.com"

Response codes can be
 * `200 OK` - If all the URLs requested were successfully fetched and returned
 * `206 Partial Content` - If some of the URLs were successfully fetched and returned
 * `204 No Content` - If all the URLs requested failed or timed out

Notes
-----

There's been no thought into security, and whatever protocol scheme's are supported by Go
are valid.


[made by jak](http://jak.io) - no license or copyright, have fun
