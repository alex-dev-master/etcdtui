package main

import "github.com/alexandr/etcdtui/internal/app/layouts"

func main() {
	m := layouts.NewManager()
	m.Render()
}
