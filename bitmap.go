package main

type BitmapRenderer struct{}

func NewBitmapRenderer() *BitmapRenderer {
	return &BitmapRenderer{}
}

func (r *BitmapRenderer) RenderToImage(*Diagram) *struct{} {
	return nil
}
