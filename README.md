## cmap - A thread-safe concurrent map

#### Install
```bash
go get github.com/s1m0n21/cmap
```

#### How to use

```go
package main

import (
	"fmt"

	"github.com/s1m0n21/cmap"
)

type key struct {
	Content int // The key must have at least one public field
}

type value struct {
	Content int
}

func main() {
	m := cmap.New(cmap.DefaultShard)

	k := key{Content: 1}
	v := value{Content: 2}

	if err := m.Set(k, v); err != nil {
		fmt.Printf("SET ERROR: %s", err)
		return
	}

	r, has, err := m.Get(k)
	if err != nil {
		fmt.Printf("GET ERROR: %s", err)
		return
	}

	fmt.Printf("has: %v", has)
	fmt.Printf("r: %+v", r)
}
```