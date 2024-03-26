package display

import (
	"context"
	"time"

	"gocv.io/x/gocv"
)

const (
	playDelayMs        = int(1.0 / 20.0 * 1_000)
	fastForwardDelayMs = int(1)
	slowMotionDelayMs  = int(1.0 / 5.0 * 1_000)
)

func CreateAndHandleGoCVWindow(ctx context.Context, cancel context.CancelFunc, mats chan gocv.Mat) error {
	defer func() {
		cancel()
	}()

	if mats == nil {
		<-ctx.Done()
		return nil
	}

	window := gocv.NewWindow("cameranator")
	defer func() {
		_ = window.Close()
	}()

	window.SetWindowProperty(gocv.WindowPropertyAspectRatio, gocv.WindowKeepRatio)
	window.MoveWindow(100, 100)

	go func() {
		<-ctx.Done()
		_ = window.Close()
	}()

	delay := playDelayMs
	pause := false

	for {
		select {
		case <-ctx.Done():
			return nil
		case mat := <-mats:
			a := time.Now().UnixMilli()

			window.IMShow(mat)
			_ = mat.Close()

			for {
				select {
				case <-ctx.Done():
					return nil
				default:
				}

				b := time.Now().UnixMilli()

				adjustedDelay := delay - int(b-a)
				if adjustedDelay < 1 {
					adjustedDelay = 1
				}

				key := window.WaitKey(adjustedDelay)

				switch key {
				case 27: // escape
					return nil

				case 32: // space
					pause = !pause

				case 0: // up
					delay = fastForwardDelayMs

				case 1: // down
					delay = slowMotionDelayMs

				case 3: // right
					delay = playDelayMs

				default:
				}

				if pause {
					continue
				}

				break
			}
		}
	}
}
