# Brain Dead
A simple Go compiler for the [Brainf*ck](https://en.wikipedia.org/wiki/Brainfuck) esoteric language.

## ⚙️ Installation
The easiest way to use this package is to install it with `go install github.com/iamscottcab/braindead`

Alternatively you can clone the repository and build it locally with `go build compiler.go`.

## ⚡️ Quick Start
1. Write some [Brainf*ck](https://en.wikipedia.org/wiki/Brainfuck) and save it to a file
2. Depending on your install method do one of the following:

**Go Install:** `braindead -in="/path/to/file`

**Local Clone & Build:** `` `./compiler.exe -in="/path/to/file`

**Local Clone:** `` `go run compiler.go -in="/path/to/file`

## 📝Configuration
While the `in` flag is required you can also optionally set the output file path and the default number of memory locations that your BF program needs to run.

**File Output:** `-out="/path/to/file"`

**Memory Size:** `-mem=1000`

## 🧠 Compiler Semantics
### Parsing
The compiler will ignore all whitespace in a program except the space literal (i.e. `' '`) when it is in a a "comment block". Take the following BF program.

```
This is a BF Program!
+          +
```

The compiler would generate the following comment in Go

```
// This is a BF Program!
```

But would interpret the following line as follows

```
++
```
### Input
The carriage return and new line characters `\r` and `\n` respectively denote the "end of input". If either is found then a zero result is returned.

Take the BF program as described in `examples/cat.bf`.

```
>,[>,]<[<]>[.>]
```

The program will continue to accept input until either `\r` or `\n` is found.

```
> This is a test // Hit Enter for new line
$ This is a test

```

This means that programs with specific input counts may suffer from odd behaviour when the input string is too short. Take the following program which expects three characters and then prints them out.

```
,.,.,.
```
Due to input semantics the program will run in one of two ways

If all input is provided immediately
```
> abc
$ abc
```
Then the program will read each input accordingly and return the string abc.

However if the user wishes to input one character at a time they would get the following behavior
```
> a // Enter
$ a
// Program waits for more input
> b // Enter
$ b
```
The first instance of the new line consumes the second input, this is by design. For programs that have fixed input values it is best to provide them in one string rather than across multiple lines.

## 💖 Thanks
If you've gotten this far, or you've enjoyed this repo and want to say thanks you can do that in the following ways:
- Add a [GitHub Star](https://github.com/iamscottcab/unity-source-generators) to the project.