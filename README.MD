![Donkey Image](img/donkey.png)

Writing an interpreter (and eventually compiler) closely based on the Monkey programming language, as described in the book [Writing An Interpreter In Go](https://interpreterbook.com/).


install the donkey repl with:
```bash
go install github.com/OJOMB/donkey/cmd/donkeyrepl@latest

donkeyrepl
```


```donkey
                          /\          /\
                         ( \\        // )
                          \ \\      // /
                           \_\\||||//_/
                            \/ _  _ \
                           \/|(O)(O)|
                          \/ |      |
      ___________________\/  \      /
     //                //     |____|
    //                ||     /      \
   //|                \|     \ 0  0 /
  // \       )         V    / \____/
 //   \     /        (     /
""     \   /_________|  |_/
       /  /\   /     |  ||
      /  / /  /      \  ||
      | |  | |        | ||
      | |  | |        | ||
      |_|  |_|        |_||
       \_\  \_\        \_\\
Welcome to the Donkey programming language!
>>
```

variable binding:
```donkey
var foo = 5;
var bar = "hello world";
var baz = true;
```

control flow:
```donkey
if (x < y) {
    return x;
} elif (x > y) {
    return y;
} else {
    return 0;
}
```

function literals and calls:
```donkey
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

closures:
```donkey
fn makeAdder(x) {
    return fn(y) { return x + y; };
}

var add5 = makeAdder(5);
add5(3); // returns 8
```

loops:
```donkey
var i = 0;
while (i < 10) {
    i = i + 1;
}

print(i); // prints 10

for (var j = 0; j < 10; j++) {
    print(j); // prints 0 through 9
}
```