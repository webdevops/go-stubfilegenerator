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
)

type stubGenerator struct {
	TemplateVariables map[string]string
	Image struct {
		NumColors int
		Quality int
		Width int
		Height int
	}
	Overwrite bool
}

func StubGenerator() stubGenerator {
	stubGen := stubGenerator{}
	stubGen.TemplateVariables = map[string]string{}
	stubGen.Image.NumColors = 16
	stubGen.Image.Quality = 70
	stubGen.Image.Width = 100
	stubGen.Image.Height = 100
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
}

func (config stubGenerator) generateStubAuto(path string) bool {
	assetFile := fmt.Sprintf("res/auto/stub%s", filepath.Ext(path))

	data, err := Asset(assetFile)
	if err == nil {
		f := createFile(path)
		defer f.Close()
		f.Write(data)
		f.Sync()
		return true
	}

	return false
}

func (config stubGenerator) generateStubImage(path string) bool {
	config.init(path)
	
	fileExt := filepath.Ext(path)
	switch fileExt {
	case ".gif":
		options := gif.Options{
			NumColors: config.Image.NumColors,
		}

		f := createFile(path)
		defer f.Close()
		img := image.NewRGBA(image.Rect(0, 0, config.Image.Width, config.Image.Height))
		gif.Encode(f, img, &options)
		return true

	case ".png":
		f := createFile(path)
		defer f.Close()
		img := image.NewRGBA(image.Rect(0, 0, config.Image.Width, config.Image.Height))
		png.Encode(f, img)
		return true

	case ".jpg":
		fallthrough
	case ".jpeg":
		options := jpeg.Options{
			Quality: config.Image.Quality,
		}

		f := createFile(path)
		defer f.Close()
		img := image.NewRGBA(image.Rect(0, 0, config.Image.Width, config.Image.Height))
		jpeg.Encode(f, img, &options)
		return true
	}

	return false
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

		f := createFile(path)
		defer f.Close()
		f.WriteString(content)
		f.Sync()

		return true
	}

	return false
}

func createFile(path string) *os.File {
	directory := filepath.Dir(path)
	os.MkdirAll(directory, os.ModePerm)
	
	f, err := os.Create(path)
	check(err)
	return f
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
