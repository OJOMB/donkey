![Donkey Image](img/donkey.png)

Writing an interpreter (and eventually compiler) closely based on the Monkey programming language, as described in the book [Writing An Interpreter In Go](https://interpreterbook.com/).

variable binding:
```monkey
var foo = 5;
var bar = "hello world";
var baz = true;
```

control flow:
```monkey
if (x < y) {
    return x;
} elif (x > y) {
    return y;
} else {
    return 0;
}
```

function literals and calls:
```monkey
var add = fn(x, y) { return x + y; };
add(2, 3); // returns 5

fn factorial(n) {
    if (n == 0) {
        return 1;
    } else {
        return n * factorial(n - 1);
    }
}

factorial(5); // returns 120
```