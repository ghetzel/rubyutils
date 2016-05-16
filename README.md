# rubyutils
A collection of utilities for performing Ruby syntax generation in Golang

## Introduction

This package introduces utilities for parsing and generating syntactically-correct Ruby from within Golang.  The main reasoning behind this is that there are numerous projects whose configuration is not a markup language (like YAML or XML) or structured-data format (such as JSON), but rather are simply executable Ruby code.  This package aims to make generating such files easier, especially with respect to outputting complex Golang data types (i.e.: slices, structs, nested maps) as valid Ruby.

## Installation and Usage

To install, run:

```
$ go get github.com/ghetzel/rubyutils
```

And import using:

```
import "github.com/ghetzel/rubyutils"
```

### Basic Usage Example

```go
package main

import (
    "fmt"
    "github.com/ghetzel/rubyutils/encoding/ruby"
)

func main() {
    myData := map[string]interface{}{
        `name`: `Test`,
        `count`: 4,
        `enabled`: true,
        `items`: []string{ `foo`, `bar`, `baz`, `qux` },
    }
    
    if data, err := ruby.MarshalIndent(myData, ``, `  `); err == nil {
        fmt.Println(string(data[:]))
    }
}
```

The above program would output:

```
{
  'count' => 4,
  'enabled' => true,
  'items' => [
    'foo',
    'bar',
    'baz',
    'qux'
],
  'name' => 'Test'
}
```
