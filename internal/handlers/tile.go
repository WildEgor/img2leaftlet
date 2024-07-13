package handlers

import (
	"fmt"
	"github.com/nao1215/imaging"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"strconv"
	"sync"
)

var _ Handler = (*TileHandler)(nil)

// TileHandler ...
type TileHandler struct {
	maxWorkers int
}

// NewTileHandler ...
func NewTileHandler() *TileHandler {
	return &TileHandler{
		4,
	}
}

// Handle ...
func (h *TileHandler) Handle(args *Args) error {
	rgbimage := imaging.Clone(args.Image)
	bounds := rgbimage.Bounds()

	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	mdim := math.Max(float64(width), float64(height))

	scaleMaxF := math.Log2(float64(mdim)/float64(args.Size)) + 1
	scaleMax := int(scaleMaxF)

	cTileSize := args.Size
	rescaleSize := args.Size

	for scale := scaleMax; scale >= 1; scale-- {
		fmt.Printf("tile size: %v\n", cTileSize)
		h.makeImageTiles(args.Output, scale, cTileSize, rescaleSize, rgbimage)
		cTileSize *= 2
	}

	return nil
}

// makeImageTiles ...
func (h *TileHandler) makeImageTiles(basePath string, scale int, tileSize int, rescaleSize int, rgbimage *image.NRGBA) {
	fmt.Println(scale, tileSize)

	bounds := rgbimage.Bounds()

	subPath := basePath + "/" + strconv.Itoa(scale)
	//nolint:gosec
	if err := os.MkdirAll(subPath, 0777); err != nil {
		return
	}

	var wg sync.WaitGroup
	jobChan := make(chan *tileJob)

	wg.Add(h.maxWorkers)
	for w := 0; w < h.maxWorkers; w++ {
		go func() {
			defer wg.Done()
			for job := range jobChan {
				processTile(job)
			}
		}()
	}

	for cx := bounds.Min.X; cx < bounds.Max.X; cx += tileSize {
		//nolint:gosec
		if err := os.MkdirAll(subPath+"/"+strconv.Itoa(cx/tileSize), 0777); err != nil {
			return
		}

		for cy := bounds.Min.Y; cy < bounds.Max.Y; cy += tileSize {
			job := &tileJob{
				cx:          cx,
				cy:          cy,
				tileSize:    tileSize,
				rescaleSize: rescaleSize,
				subPath:     subPath,
				rgbimage:    rgbimage,
			}
			jobChan <- job
		}
	}

	close(jobChan)
	wg.Wait()
}

// tileJob ...
type tileJob struct {
	cx, cy, tileSize, rescaleSize int
	subPath                       string
	rgbimage                      *image.NRGBA
}

// processTile ...
func processTile(job *tileJob) {
	subimage := job.rgbimage.SubImage(image.Rectangle{
		Min: image.Point{
			X: job.cx,
			Y: job.cy,
		},
		Max: image.Point{
			X: job.cx + job.tileSize,
			Y: job.cy + job.tileSize,
		},
	}).(*image.NRGBA)

	subbounds := subimage.Bounds()
	xDelta := subbounds.Max.X - subbounds.Min.X
	yDelta := subbounds.Max.Y - subbounds.Min.Y

	if xDelta < job.tileSize || yDelta < job.tileSize {
		newSubImage := image.NewNRGBA(image.Rectangle{
			Min: image.Point{},
			Max: image.Point{
				X: job.tileSize,
				Y: job.tileSize,
			},
		})
		draw.Draw(newSubImage, image.Rectangle{
			Min: image.Point{},
			Max: image.Point{
				X: job.tileSize,
				Y: job.tileSize,
			},
		}, subimage, subimage.Bounds().Min, draw.Src)
		subimage = newSubImage
	}

	if job.tileSize != job.rescaleSize {
		subimage = imaging.Resize(subimage, job.rescaleSize, job.rescaleSize, imaging.Lanczos)
	}

	subfile, err := os.Create(job.subPath + "/" + strconv.Itoa(job.cx/job.tileSize) + "/" + strconv.Itoa(job.cy/job.tileSize) + ".png")
	if err != nil {
		fmt.Println("Failed to create tile file:", err)
		return
	}
	defer subfile.Close() //nolint:all

	if err := png.Encode(subfile, subimage); err != nil {
		fmt.Println("Failed to encode tile image:", err)
		return
	}
}
