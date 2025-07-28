package rgb

import (
	"log"
	"math"
	"time"
	"golang.org/x/exp/slices"
)

// Spinner will run RGB function
func (r *ActiveRGB) Spinner(startTime *time.Time) {
	//Time elapsed in milliseconds
	elapsed := time.Since(*startTime).Milliseconds()
	//Elapsed time divided by ( speed multiplied by 100
	progress := math.Mod(float64(elapsed)/(r.RgbModeSpeed*1000), 1.0)
	log.Println("Channels")
	log.Println(r.LightChannels)
	if progress >= 1.0 {
		*startTime = time.Now() // Reset startTime to the current time
		elapsed = 0             // Reset elapsed time
		progress = 0            // Reset progress
	}
	activeLEDs := []int{0, 0, 0, 0}
	activeLEDs[0] = int(progress * float64(r.LightChannels))
	activeLEDs[1] = (activeLEDs[0] + 4) % r.LightChannels
	activeLEDs[2] = (activeLEDs[0] + 1) % r.LightChannels
	activeLEDs[3] = (activeLEDs[1] + 1) % r.LightChannels
	buf := map[int][]byte{}
	for j := 0; j < r.LightChannels; j++ {
		log.Println("J value")
		log.Println(j)
		if len(r.Buffer) > 0 {
			// If colour values are already in the buffer...
			isActive := slices.Contains(activeLEDs, j)
			if isActive {
				t := float64(activeLEDs[0]) / float64(40)
				colors := interpolateColors(r.RGBStartColor, r.RGBStartColor, t, r.RGBBrightness)
				r.Buffer[j] = byte(colors.Red)
				r.Buffer[j+r.ColorOffset] = byte(colors.Green)
				r.Buffer[j+(r.ColorOffset*2)] = byte(colors.Blue)
			} else {
				r.Buffer[j] = 255
				r.Buffer[j+r.ColorOffset] = 255
				r.Buffer[j+(r.ColorOffset*2)] = 0
			}
		} else {
			// Fill the buffer with colour values
			isActive := slices.Contains(activeLEDs, j)
			if isActive {
				t := float64(activeLEDs[0]) / float64(40)
				colors := interpolateColors(r.RGBStartColor, r.RGBStartColor, t, r.RGBBrightness)
				buf[j] = []byte{
					byte(colors.Red),
					byte(colors.Green),
					byte(colors.Blue),
				}

				if r.IsAIO && r.HasLCD {
					if j > 15 && j < 20 {
						buf[j] = []byte{0, 0, 0}
					}
				}
			} else {
				t := float64(activeLEDs[0]) / float64(40)
				colors := interpolateColors(r.RGBEndColor, r.RGBEndColor, t, r.RGBBrightness)
				buf[j] = []byte{
					byte(colors.Red),
					byte(colors.Green),
					byte(colors.Blue),
				}

			}
		}
	}
	// Raw colors
	r.Raw = buf

	if r.Inverted {
		r.Output = SetColorInverted(buf)
	} else {
		r.Output = SetColor(buf)
	}
}

