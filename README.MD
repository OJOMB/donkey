![Donkey Image](img/donkey.png)

An interpreter (and eventually compiler) that is closely based on the Monkey programming language, as described in the book [Writing An Interpreter In Go](https://interpreterbook.com/).

Donkey expands on Monkey with additional features such as:
* for/while loops
* function statements
* variable rebinding
* bitwise operators
* statement block scoping
* increment/decrement suffix operators
* and more!

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

variable binding and rebinding:
```donkey
var foo = 5;
foo = 10; // rebinds foo to 10

var bar = "hello world";
bar = "goodbye world"; // rebinds bar to "goodbye world"

var baz = true;
baz = false; // I'm sure you can guess what this does
```

lists:
```donkey
var myList = [1, 2, 3, 4, 5];
print(myList[0]); // prints 1
print(myList[1]); // prints 2
print(myList[2]); // prints 3

// lisp-style list operations:
cdr(myList); // returns [2, 3, 4, 5]
cons(0, myList); // returns [0, 1, 2, 3, 4, 5]
car(myList); // returns 1
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

function literals, function statements, and calls:
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

bitwise operators:
```donkey
var a = 5; // 0101 in binary
var b = 3; // 0011 in binary
print(a & b); // prints 1 (0001 in binary)

print(a | b); // prints 7 (0111 in binary)
print(a ^ b); // prints 6 (0110 in binary)
print(a << 1); // prints 10 (1010 in binary)
print(a >> 1); // prints 2 (0010 in binary)
```