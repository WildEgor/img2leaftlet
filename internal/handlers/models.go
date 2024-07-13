// Package handlers ...
package handlers

import "image"

// Args ...
type Args struct {
	Output string
	Size   int

	Image image.Image
}

// Handler ...
type Handler interface {
	Handle(args *Args) error
}
