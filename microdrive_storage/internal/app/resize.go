package app

import (
	"bytes"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"os"
	"strconv"
	sync "sync"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
)

func resizeAndSave(paths []OriginalImage) []string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	uploadPaths := make([]string, 0)
	for i := range paths {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			uploadPath := resizeImage(paths[i])
			if len(uploadPath) != 0 {
				mu.Lock()
				uploadPaths = append(uploadPaths, uploadPath...)
				mu.Unlock()
			} else {
				mu.Lock()
				uploadPaths = append(uploadPaths, paths[i].Path)
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()
	return uploadPaths
}

// наша основная функция на этом этапе обработки
func resizeImage(path OriginalImage) []string {
	var mu sync.Mutex
	// проверка что длина и ширина изображения имеются
	for i := range path.Lenght {
		if path.Lenght[i] == 0 || path.Width[i] == 0 {
			return []string{}
		}
	}
	//проверяем что у нас одинаковое количество
	// параметров длины и ширины картинки для ее создания
	if len(path.Lenght) != len(path.Width) {
		return []string{}
	}
	//результирующий слайс путей наших файлов
	var uploadPaths []string
	switch path.Format {
	case "png":
		for i := range path.Lenght {
			//открытие файла уже сжатого и с ватермаркой
			imgIn, err := os.Open(path.Path)
			if err != nil {
				slog.Error("failed to open PNG file", "error", err.Error())
				return []string{}
			}
			//его считывание
			imgPng, err := png.Decode(imgIn)
			if err != nil {
				slog.Error("failed to decode PNG", "error", err.Error())
				return []string{}
			}
			// и закрытие
			err = imgIn.Close()
			if err != nil {
				slog.Error("failed to close PNG file", "error", err.Error())
				return []string{}
			}
			// генерация нового изображения с измененным размером
			imgPng = resize.Resize(uint(path.Lenght[i]), uint(path.Width[i]), imgPng, resize.Bilinear)
			upPath := path.Folder + path.UUID + "_" + strconv.FormatUint(uint64(path.Lenght[i]), 10) + "x" + strconv.FormatUint(uint64(path.Width[i]), 10) + "." + path.Format
			buf := new(bytes.Buffer)
			// его запись в нужный нам формат
			err = png.Encode(buf, imgPng)
			if err != nil {
				slog.Error("failed to encode PNG", "error", err.Error())
				return []string{}
			}
			imgSave := buf.Bytes()
			//сохранение нового изображения по сгенерированному пути
			err = os.WriteFile(upPath, imgSave, 0666)
			if err != nil {
				slog.Error("failed to save PNG file", "error", err.Error())
				return []string{}
			}
			// добавление полученного пути к общему слайсу
			mu.Lock()
			uploadPaths = append(uploadPaths, upPath)
			mu.Unlock()
		}
		// по аналогии с кейсом png
	case "jpg", "jpeg":
		for i := range path.Lenght {
			imgIn, err := os.Open(path.Path)
			if err != nil {
				slog.Error("failed to open JPEG file", "error", err.Error())
				return []string{}
			}
			imgJpeg, err := jpeg.Decode(imgIn)
			if err != nil {
				slog.Error("failed to decode JPEG", "error", err.Error())
				return []string{}
			}
			err = imgIn.Close()
			if err != nil {
				slog.Error("failed to close JPEG file", "error", err.Error())
				return []string{}
			}
			imgJpeg = resize.Resize(uint(path.Lenght[i]), uint(path.Width[i]), imgJpeg, resize.Bilinear)
			upPath := path.Folder + path.UUID + "_" + strconv.FormatUint(uint64(path.Lenght[i]), 10) + "x" + strconv.FormatUint(uint64(path.Width[i]), 10) + "." + path.Format
			buf := new(bytes.Buffer)
			err = jpeg.Encode(buf, imgJpeg, &jpeg.Options{Quality: 100})
			if err != nil {
				slog.Error("failed to encode JPEG", "error", err.Error())
				return []string{}
			}
			imgSave := buf.Bytes()
			err = os.WriteFile(upPath, imgSave, 0666)
			if err != nil {
				slog.Error("failed to save JPEG file", "error", err.Error())
				return []string{}
			}
			mu.Lock()
			uploadPaths = append(uploadPaths, upPath)
			mu.Unlock()
		}
		// по аналогии с кейсом png
	case "webp":
		for i := range path.Lenght {
			imgIn, err := os.Open(path.Path)
			if err != nil {
				slog.Error("failed to open WEBP file", "error", err.Error())
				return []string{}
			}
			// Decode webp
			mg, err := io.ReadAll(imgIn)
			if err != nil {
				slog.Error("failed to read WEBP file", "error", err.Error())
				return []string{}
			}
			m, err := webp.DecodeRGB(mg)
			if err != nil {
				slog.Error("failed to decode WEBP", "error", err.Error())
				return []string{}
			}
			err = imgIn.Close()
			if err != nil {
				slog.Error("failed to close WEBP file", "error", err.Error())
				return []string{}
			}
			imgWebp := resize.Resize(uint(path.Lenght[i]), uint(path.Width[i]), m, resize.Bilinear)
			upPath := path.Folder + path.UUID + "_" + strconv.FormatUint(uint64(path.Lenght[i]), 10) + "x" + strconv.FormatUint(uint64(path.Width[i]), 10) + "." + path.Format
			imgSave, err := webp.EncodeRGB(imgWebp, 100)
			if err != nil {
				slog.Error("failed to encode WEBP", "error", err.Error())
				return []string{}
			}
			err = os.WriteFile(upPath, imgSave, 0666)
			if err != nil {
				slog.Error("failed to save WEBP file", "error", err.Error())
				return []string{}
			}
			mu.Lock()
			uploadPaths = append(uploadPaths, upPath)
			mu.Unlock()
		}
	default:
		slog.Error("unknown file format, please provide a file with extension png,jpg,webp")
	}
	return uploadPaths
}
