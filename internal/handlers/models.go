package handlers

import "image"

type Args struct {
	Output string
	Size   int

	Image image.Image
}

type Handler interface {
	Handle(args *Args) error
}
