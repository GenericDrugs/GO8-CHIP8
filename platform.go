package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type Platform struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
}

func NewPlatform(title string, windowWidth, windowHeight, textureWidth, textureHeight int32) (*Platform, error) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, err
	}

	p := &Platform{}

	var err error
	p.window, err = sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		windowWidth,
		windowHeight,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		sdl.Quit()
		return nil, err
	}

	p.renderer, err = sdl.CreateRenderer(p.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		p.Destroy()
		return nil, err
	}

	p.texture, err = p.renderer.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_STREAMING,
		textureWidth,
		textureHeight,
	)
	if err != nil {
		p.Destroy()
		return nil, err
	}

	return p, nil
}

func (p *Platform) Destroy() {
	if p.texture != nil {
		p.texture.Destroy()
	}
	if p.renderer != nil {
		p.renderer.Destroy()
	}
	if p.window != nil {
		p.window.Destroy()
	}
	sdl.Quit()
}

func (p *Platform) Update(buffer []uint32, pitch int) error {
	if len(buffer) == 0 {
		return fmt.Errorf("video buffer is empty")
	}

	cBuffer := C.malloc(C.size_t(len(buffer) * 4))
	if cBuffer == nil {
		return fmt.Errorf("could not allocate SDL pixel buffer")
	}
	defer C.free(cBuffer)

	// Copy Go pixel data into C-allocated memory.
	goPixels := unsafe.Slice((*uint32)(cBuffer), len(buffer))
	copy(goPixels, buffer)

	if err := p.texture.Update(nil, cBuffer, pitch); err != nil {
		return err
	}
	if err := p.renderer.Clear(); err != nil {
		return err
	}
	if err := p.renderer.Copy(p.texture, nil, nil); err != nil {
		return err
	}

	p.renderer.Present()
	return nil
}

var chip8KeyMap = map[sdl.Keycode]uint8{
	sdl.K_x: 0x0,
	sdl.K_1: 0x1, sdl.K_2: 0x2, sdl.K_3: 0x3,
	sdl.K_q: 0x4, sdl.K_w: 0x5, sdl.K_e: 0x6,
	sdl.K_a: 0x7, sdl.K_s: 0x8, sdl.K_d: 0x9,
	sdl.K_z: 0xA, sdl.K_c: 0xB,
	sdl.K_4: 0xC, sdl.K_r: 0xD, sdl.K_f: 0xE, sdl.K_v: 0xF,
}

func (p *Platform) ProcessInput(keys *[16]uint8) bool {
	quit := false

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			quit = true

		case *sdl.KeyboardEvent:
			key := e.Keysym.Sym
			if e.Type == sdl.KEYDOWN && key == sdl.K_ESCAPE {
				quit = true
			}
			if index, ok := chip8KeyMap[key]; ok {
				if e.Type == sdl.KEYDOWN {
					keys[index] = 1
				} else if e.Type == sdl.KEYUP {
					keys[index] = 0
				}
			}
		}
	}

	return quit
}
