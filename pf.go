package pf

import (
	"sync"
)

// Perform an operation on a single ARGB pixel
type PixelFunction func(v uint32) uint32

// Combine two functions to a single PixelFunction.
// The functions are applied in the same order as the arguments.
func Combine(a, b PixelFunction) PixelFunction {
	return func(v uint32) uint32 {
		return b(a(v))
	}
}

// Combine three functions to a single PixelFunction.
// The functions are applied in the same order as the arguments.
func Combine3(a, b, c PixelFunction) PixelFunction {
	return Combine(a, Combine(b, c))
}

// Map a PixelFunction over every pixel (uint32 ARGB value)
func Map(cores int, f PixelFunction, pixels []uint32) {
	// Map a pixel function over every pixel, concurrently

	var wg sync.WaitGroup

	iLength := int32(len(pixels))

	iStep := iLength / int32(cores)
	iConcurrentlyDone := int32(cores) * iStep
	var iDone int32

	// Apply partialMap for each of the partitions
	if iStep < iLength {
		for i := int32(0); i < iConcurrentlyDone; i += iStep {
			wg.Add(1)
			go partialMap(&wg, f, pixels, i, i+iStep)
		}
		iDone = iConcurrentlyDone
	}

	if iDone == iLength {
		// No leftover pixels
		return
	}

	// Apply partialMap to the final leftover pixels
	wg.Add(1)
	go partialMap(&wg, f, pixels, iDone, iLength)

	wg.Wait()
}
