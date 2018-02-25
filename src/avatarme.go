package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	//"strings"
	"image"
	"image/color"
	"image/png"
)

const (
	w = 256
)

func main() {
	fmt.Println(os.Args[1])

	bArray := sha256.Sum256([]byte(os.Args[1]))

	fmt.Println(bArray)

	img := image.NewNRGBA(image.Rect(0, 0, w, w))
	bLen := len(bArray)
	loop := w / bLen
	for i := 0; i < loop; i++ {
		for j := 0; j < loop; j++ {
			wStart := i * bLen
			hStart := j * bLen
			for m := 0; m < bLen; m++ {
				for n := 0; n < bLen; n++ {
					v := bArray[(i*loop+j)%32]
					img.Set(wStart+m, hStart+n, color.NRGBA{v, v, v, 255})
				}
			}
		}
	}

	f, _ := os.Create("avatarme.png")
	defer f.Close()
	png.Encode(f, img)
}
