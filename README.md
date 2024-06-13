# Go Assert

Go Assert provides assertions and a cleanup tool for contract based programming.

## Description

Assertions are not a replacement for error handling. An assert is typically used
to verify that the state of the program is valid at all times. If an unexpected invalid
state occurs, the assert will catch the error and crash with a stack trace
containing the invalid state.

Additionally, asserts can be removed prior to release for performance. Arguably,
asserts are no longer necessary once the program has been thoroughly tested for
invalid state.

This library also provides 'dass' (short for delete assertions), a tool which
can safely remove assertions recursively from a given root folder.

## Getting Started

### Dependencies

* Go version 1.18 or higher.

### Installing

```zsh
go get github.com/jdavasligil/go-assert
```

### Executing program

By default, dass will only search locally within the current working directory
and it will NOT preserve backups.

```zsh
go run github.com/jdavasligil/go-assert/cmd/dass
```


## Help
For a list of options, run with the `-h` flag.

```
go run github.com/jdavasligil/go-assert/cmd/dass -h
```

## Authors

Jaedin Davasligil
[contact](jdavasligil.swimming625@slmails.com)

## Version History

* 1.0.0 

## License
The MIT License (MIT)

Copyright (c) 2024 J. Davasligil

See LICENSE.md for details.

## Acknowledgments

* [awesome-readme](https://github.com/matiassingers/awesome-readme)
* [pitchfork](https://api.csswg.org/bikeshed/?force=1&url=https://raw.githubusercontent.com/vector-of-bool/pitchfork/develop/data/spec.bs)
* [project-layout](https://github.com/golang-standards/project-layout)
