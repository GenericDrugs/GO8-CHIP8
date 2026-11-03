package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <Scale> <Delay> <ROM>\n", os.Args[0])
		os.Exit(1)
	}

	videoScale, err := strconv.Atoi(os.Args[1])
	if err != nil || videoScale <= 0 {
		fmt.Fprintln(os.Stderr, "Scale must be a positive integer")
		os.Exit(1)
	}

	cycleDelay, err := strconv.Atoi(os.Args[2])
	if err != nil || cycleDelay < 0 {
		fmt.Fprintln(os.Stderr, "Delay must be a non-negative integer")
		os.Exit(1)
	}

	romFilename := os.Args[3]

	platform, err := NewPlatform(
		"CHIP-8 Emulator",
		int32(VIDEO_WIDTH*videoScale),
		int32(VIDEO_HEIGHT*videoScale),
		int32(VIDEO_WIDTH),
		int32(VIDEO_HEIGHT),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer platform.Destroy()

	c8 := newChip8()
	if err := c8.LoadROM(romFilename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	videoPitch := 4 * VIDEO_WIDTH // sizeof(uint32) * width
	lastCycleTime := time.Now()
	var quit bool

	for !quit {
		quit = platform.ProcessInput(&c8.keypad)

		if time.Since(lastCycleTime) > time.Duration(cycleDelay)*time.Millisecond {
			lastCycleTime = time.Now()

			c8.Cycle()

			if err := platform.Update(c8.video[:], videoPitch); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}
	}
}
