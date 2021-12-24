# takeout-browser

## Description

Takeout-Browser is a Go web-based application to browse the contents of a [Google Takeout](https://takeout.google.com/) export. For now, the following features are fully implemented:

- None

The following features are being worked on:

- Google Photos gallery

The following features are in the backlog:

- Google Drive browser
- Hangouts viewer
- Keep viewer

## Usage

This app requires a Google Takeout archive containing a Google Photos extract. You can follow [this tutorial](https://support.google.com/accounts/answer/9666875).

Then, extract the archive somewhere and specify the path to the `Takeout` folder when running the app:

```
$ ./takeout-browser -folder /path/to/Takeout
```

Then, head on to http://127.0.0.1:8080 to browse your pictures.

## Building

Requires Go 1.17.2 or newer.

```
$ go build -o takeout-browser cmd/main.go
```