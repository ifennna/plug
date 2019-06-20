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
 
This is Plug's WebAssembly implementation that allows Plug code to be 
run in the browser.

To run code, serve the `web` folder.