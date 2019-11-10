# Dashboard

> A simple and extensible dashboard to display data that matters to you

**WARNING: THIS IS ALPHA STAGE SOFTWARE**

Dashboard is designed to allow developers to build widgets by writing plug-ins. Dashboard supports plug-ins written in JSON out-of-the-box. Support for other technologies such as TOML, YAML etc will be added in the future. The goal is to allow developers to develop their own custom plug-ins which can be shared among the members of the Dashboard community. The project is actively being developed

![](https://i.imgur.com/r4Dla3J.png)


## Getting Started

* [TODO] Hosting
* Push data to `/push` endpoint

## Installation

* [TODO] Deploy Dashboard server

## Features

* Push support
* Pull support
* Extensible via plug-ins
* Shareable plug-ins
* [TODO] Map response values to display in Dashboard
* [TODO] Category support
* [TODO] Response Hooks

## Development

For convenience, a `Makefile` is setup. From the app root, running the `make` command will build and run the go program.

If you want to run the commands manually, run `go build` from the app root. Once everything is built, run `PORT=300 ./dashboard` from the app root. You should be able to access Dashboard at [http://localhost:3000](http://localhost:3000). PLug-ins can be written and tested out locally.

## Tests

Run `go test`

## Technology

Dashboard is built on the following stack:

* Back-end: Go
* Front-end: React [TODO]

## Plug-ins

A Dashboard Plug-in is a simple config file stores details on how to access a particular endpoint and which parts of the response are required. This can be written by anyone to access any data that they'd like to display using Dashboard.

Following are the supported fields

- `id`: unique identifier to identify the plug-in
- `name`: display name of the component in the UI
- `url`: endpoint that dashboard needs to poll
- `method` (optional. default is `GET`): restful verb that needs to be used to query the endpoint [GET/POST]
- `headers` (optional): headers that need to be passed in for the request
- `interval` (optional. default is `10` seconds): polling interval

Here's an example:

```json
{
  "id": "unique-id",
  "name": "Server Uptime",
  "method": "GET",
  "headers": {},
  "url": "http://localhost:3002/data",
  "interval": "5",
}

```

## Future

* Category/Page Support

Group components and display them as separate pages. This enables a single Dashboard server to handle multiple subscriptions with the convenience of displaying the components based on categories.

* Filterable Subscription values

Currently, Dashboard expects a certain structure in the JSON responses of the pull calls. In the future, the plug-in authors will be able to filter and map JSON values from the response.

* Hooks

Each response that Dashboard receives from either a push or a pull, will be parsed and a hook can be run, if configured.

For example, say we have a poll subscription that checks for the uptime of a particular endpoint. If Dashboard receives a response that's not a 200, it can run the configured hooks. These hooks can either be a message to post to Slack or notify the devs or just pop up a message alerting the people that need to be notified.

* Support for multiple languages

TOML and YAML (partially supported already since JSON is supported)

* Support human readable text in intervals

Currently, Dashboard only supports integer values to represent the number of seconds. Support for text such as `1m`, `1h`, `1m20s`, `5s` etc will be supported soon.

## Contributing

* [TODO] Add Guidelines here

## License

(The MIT License)

Copyright (c) 2019 Mohnish Thallavajhula <hi@iam.mt>

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
