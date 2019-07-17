# Plug

[![Build Status](https://travis-ci.org/noculture/plug.svg?branch=master)](https://travis-ci.org/noculture/plug)

Plug is a tiny C-like programming language. Plug syntax looks like 
this:

```$xslt
let five = 5;
let ten = 10;

let add = func(x, y) {
    x + y;
};

let result = add(five, ten);
```
 
This is basically a toy project for my personal edification. It's very much 
in its infancy.

To run the compile the project, make sure you have [Go](https://golang.org/dl/) installed. Clone the project to your `$GOPATH` and run `go build` in the project folder. 

If you have the binary already run `plug your-file.plug` or just run `plug` to start the REPL.
