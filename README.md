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
in its infancy. The only thing available is a REPL.

To run the REPL, make sure you have [Go](https://golang.org/dl/) installed. Clone the project to 
your `$GOPATH` and run `go run main.go` from the project folder. 