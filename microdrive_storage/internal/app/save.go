package app

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log/slog"
	pb "microdrive_storage/internal/gen/go"
	"os"
	sync "sync"
	"time"

	"github.com/chai2010/webp"
	"github.com/google/uuid"
)

func getTimeData() (string, string, string) {
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	y := fmt.Sprintf("%v", year)
	m2 := int(month)
	m := fmt.Sprintf("%v", m2)
	d := fmt.Sprintf("%v", day)
	return y, m, d
}

func setLevelCompressionPNG(level string) png.CompressionLevel {
	var compessInt png.CompressionLevel
	switch level {
	case "low":
		compessInt = -3
	case "medium":
		compessInt = -2
	case "max":
		compessInt = -1
	}
	return compessInt
}

func setLevelCompressionJPG(level string) int {
	var compessInt int
	switch level {
	case "low":
		compessInt = 10
	case "medium":
		compessInt = 50
	case "max":
		compessInt = 100
	}
	return compessInt
}

func setLevelCompressionWEBP(level string) float32 {
	var compessInt float32
	switch level {
	case "low":
		compessInt = 10
	case "medium":
		compessInt = 50
	case "max":
		compessInt = 100
	}
	return compessInt
}

func saveSourceFiles(images []*pb.DownloadImagesRequest) []OriginalImage {
	var wg sync.WaitGroup
	var mu sync.Mutex
	y, m, d := getTimeData()
	imagesNew := make([]OriginalImage, 0)
	for i := range images {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			extension := images[i].Info.Format
			uuid := uuid.New().String()
			newFileName := uuid + "." + extension
			newFilePath := "./download/" + y + "/" + m + "/" + d + "/" + uuid + "/" + "img/"
			//создаём папку для сохранения файлов
			err := os.MkdirAll(newFilePath, 0755)
			if err != nil {
				slog.Error("failed to create directory", "error", err.Error())
			}
			// декодируем изображение
			img, _, err := image.Decode(bytes.NewReader(images[i].Image))
			if err != nil {
				slog.Error("failed to decode image", "error", err.Error())
				return
			}
			//собираем полный путь до файла с его названием и расширением
			path := newFilePath + newFileName
			//создаём файл для последующего сохранения туда картинки
			out, _ := os.Create(path)
			defer out.Close()
			//теперь обрабатываем картинку в зависимости от полученного формата изображения
			switch extension {
			case "png":
				var enc png.Encoder
				//сжимаем в соответствии с полученными параметрами от клиента
				level := setLevelCompressionPNG(images[i].Info.Compress)
				enc.CompressionLevel = level
				//сохраняем изображение
				err = enc.Encode(out, img)
				if err != nil {
					slog.Error("failed to encode PNG", "error", err.Error())
				}
			case "jpeg", "jpg":
				var opts jpeg.Options
				level := setLevelCompressionJPG(images[i].Info.Compress)
				opts.Quality = level
				err = jpeg.Encode(out, img, &opts)
				if err != nil {
					slog.Error("failed to encode JPEG", "error", err.Error())
				}
			case "webp":
				var data []byte
				level := setLevelCompressionWEBP(images[i].Info.Compress)
				data, err = webp.EncodeRGB(img, level)
				if err != nil {
					slog.Error("failed to encode WEBP", "error", err.Error())
				}
				if err = os.WriteFile(path, data, 0666); err != nil {
					slog.Error("failed to write WEBP file", "error", err.Error())
				}
			default:
				slog.Error("unknown file format, please provide a file with extension png,jpg,webp")
			}
			imageNew := OriginalImage{
				Path:      path,
				Lenght:    images[i].Info.Height,
				Width:     images[i].Info.Width,
				Format:    images[i].Info.Format,
				Folder:    newFilePath,
				Watermark: images[i].Info.Watermark,
				UUID:      uuid,
			}
			mu.Lock()
			imagesNew = append(imagesNew, imageNew)
			mu.Unlock()
		}(i)
		wg.Wait()
	}
	return imagesNew
}
