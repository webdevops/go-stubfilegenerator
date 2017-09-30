package stubfilegenerator

import (
	"path/filepath"
	"fmt"
	"strings"
	"os"
	"image"
	"image/png"
	"image/gif"
	"image/jpeg"
	"image/color"
	"image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/font/basicfont"
)

type stubGenerator struct {
	TemplateVariables map[string]string
	Image struct {
		BackgroundColor color.RGBA
		FontColor color.RGBA
		Font *basicfont.Face
		FontLineHeight int
		NumColors int
		Quality int
		Width int
		Height int
		Text []string
	}
	Overwrite bool
}

func StubGenerator() stubGenerator {
	stubGen := stubGenerator{}
	stubGen.TemplateVariables = map[string]string{}
	stubGen.Image.BackgroundColor = color.RGBA{255, 255, 170, 255}
	stubGen.Image.FontColor = color.RGBA{0, 0, 0, 255}
	stubGen.Image.Font = inconsolata.Regular8x16
	stubGen.Image.FontLineHeight = 15
	stubGen.Image.NumColors = 256
	stubGen.Image.Quality = 80
	stubGen.Image.Width = 600
	stubGen.Image.Height = 200
	stubGen.Image.Text = []string{
		"Path: %PATH%",
	}
	stubGen.Overwrite = false
	return stubGen
}

func (config stubGenerator) init(path string) {
	if _, ok := config.TemplateVariables["PATH"]; !ok {
		config.TemplateVariables["PATH"] = path
	}
}

func (config stubGenerator) GenerateStub(path string) {

	if config.Overwrite == false {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			// file eixsts
			return
		}
	}

	if config.generateStubAuto(path) {
		return
	}

	if config.generateStubImage(path) {
		return
	}

	if config.generateStubText(path) {
		return
	}

	config.generateStubFallback(path)
}

func (config stubGenerator) generateStubAuto(path string) bool {
	assetFile := fmt.Sprintf("res/auto/stub%s", filepath.Ext(path))

	data, err := Asset(assetFile)
	if err == nil {
		createFile(path, func(f *os.File) {
			f.Write(data)
		})
		return true
	}

	return false
}

func (config stubGenerator) generateStubImage(path string) bool {
	config.init(path)
	
	fileExt := filepath.Ext(path)
	switch fileExt {
	case ".gif":
		createFile(path, func(f *os.File) {
			gif.Encode(f, config.createImage(), &gif.Options{config.Image.NumColors, nil, nil})
		})
		return true

	case ".png":
		createFile(path, func(f *os.File) {
			png.Encode(f, config.createImage())
		})
		return true

	case ".jpg":
		fallthrough
	case ".jpeg":
		options := jpeg.Options{
			Quality: config.Image.Quality,
		}

		createFile(path, func(f *os.File) {
			jpeg.Encode(f, config.createImage(), &options)
		})

		return true
	}

	return false
}

func (config stubGenerator) createImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, config.Image.Width, config.Image.Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{config.Image.BackgroundColor}, image.ZP, draw.Src)

	for num, line := range config.Image.Text {
		for key, value := range config.TemplateVariables {
			replaceKey := fmt.Sprintf("%%%s%%", key)
			line = strings.Replace(line, replaceKey, value, -1)
		}

		posX := 10
		posY := 10 + (config.Image.FontLineHeight * (num + 1))
		config.imageWriteTextLine(img, posX, posY, line)
	}

	return img
}

func (config stubGenerator) generateStubText(path string) bool {
	config.init(path)
	assetFile := fmt.Sprintf("res/templates/stub%s", filepath.Ext(path))

	data, err := Asset(assetFile)
	if err == nil {
		content := string(data)
		for key, value := range config.TemplateVariables {
			replaceKey := fmt.Sprintf("%%%s%%", key)
			content = strings.Replace(content, replaceKey, value, -1)
		}

		createFile(path, func(f *os.File) {
			f.WriteString(content)
		})

		return true
	}

	return false
}

func (config stubGenerator) generateStubFallback(path string) bool {
	createFile(path, func(f *os.File) {
		for key, value := range config.TemplateVariables {
			f.WriteString(fmt.Sprintf("%s: %s", key, value))
		}
	})

	return true
}

func createFile(path string, callback func(f *os.File)) *os.File {
	directory := filepath.Dir(path)
	os.MkdirAll(directory, os.ModePerm)
	
	f, err := os.Create(path)
	check(err)
	defer func() {
		f.Close()
		f.Sync()
	}()

	callback(f)

	return f
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (config stubGenerator) imageWriteTextLine(img *image.RGBA, x, y int, label string) {
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(config.Image.FontColor),
		Face: config.Image.Font,
		Dot:  point,
	}
	d.DrawString(label)
}
