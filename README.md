# ECE: Entity Component Event

ECE is a small, hand-made **Entity Component Event** library written in Go.

## Description

ECE is a variant of the classic Entity Component System (ECS) where the
**Systems are removed** and everything is driven by **Events**:

- **Entities** are plain integer IDs.
- **Components** are stored in a `SparseArray`, one per component type.
- **Events** are the only way to trigger behaviour. Subscribe a callback to an
  event, then emit it immediately or queue it and process it later (for example
  once per frame).

This package is mainly built for our **D&Gine** project: a free alternative to
play Dungeons & Dragons online, with no limitation, and modular enough to fit
any campaign and scenario, with many customisation options.

## How to install

Requires **Go 1.27.1** or newer.

```bash
go get github.com/<your-org>/<your-repo>
```

Then import it in your code:

```go
import "github.com/<your-org>/<your-repo>"
```

> Replace the path above with the real module path of this repository.

## Example

```go
package main

import (
	"fmt"

	"github.com/<your-org>/<your-repo>"
)

type Damage struct {
	Target int
	Amount int
}

type Health struct {
	Points int
}

func main() {
	r := ece.NewRegistry()

	healths := ece.NewSparseArray[*Health]()
	healths.Insert(1, &Health{Points: 100})

	r.SubscribeToEventByName("damage", func(r *ece.Registry, e ece.Event) any {
		d := e.(Damage)
		if h, err := healths.Get(d.Target); err == nil {
			h.Points -= d.Amount
			fmt.Printf("entity %d took %d damage, %d HP left\n", d.Target, d.Amount, h.Points)
		}
		return nil
	}, 0)

	// Handled immediately.
	r.EmitEventByName("damage", Damage{Target: 1, Amount: 10})

	// Queued, then handled frame by frame. A BreakpointEvent ends a frame.
	r.EnqueueEventByName("damage", Damage{Target: 1, Amount: 20})
	r.EnqueueEvent(&ece.BreakpointEvent{})
	r.EnqueueEventByName("damage", Damage{Target: 1, Amount: 5})

	r.ProcessEvents() // frame 1: handles the 20 damage
	r.ProcessEvents() // frame 2: handles the 5 damage
}
```

The same example, runnable as a test, is in
[`example_test.go`](example_test.go).

## How to launch documentation

The code is documented with standard Go doc comments.

In the terminal:

```bash
go doc -all .                      # whole package
```
```bash
go doc SparseArray                 # one type
```
```bash
go doc EventQueue.Process          # one method
```

In the browser, with a local pkgsite server:

```bash
go install golang.org/x/pkgsite/cmd/pkgsite@latest
pkgsite -open .
```

Run the documented example (checks its expected output):

```bash
go test -run Example -v
```

## Contributors

| Name | Role | Contact |
| ---- | ---- | ------- |
| _Your name_ | Creator / Maintainer | _link or email_ |
| @sami-hmr | Creator | sami.hamrouni@epitech.eu |
| @JulesDELTEIL | Creator | jules.delteil@epitech.eu |
