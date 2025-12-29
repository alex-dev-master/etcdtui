package main

import (
	"context"
	"fmt"

	"github.com/alexandr/etcdtui/internal/app/layouts"
)

func main() {
	m := layouts.NewManager()
	var err error
	ctx := context.Background()
	if err = m.Render(ctx); err != nil {
		fmt.Println(err)
		return
	}
}
