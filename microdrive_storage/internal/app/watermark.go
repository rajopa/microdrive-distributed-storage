package app

import (
	"os"
	sync "sync"

	"github.com/disintegration/imaging"
	"github.com/filipenevs/go-imagewatermark"
)

func watermark(paths []OriginalImage) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var err error
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}
	for i := range paths {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			watermarkS := ""

			switch paths[i].Watermark {
			case defaultWatermark:
				watermarkS = defaultWatermark
			case "":
				watermarkS = defaultWatermark
			default:
				watermarkS = paths[i].Watermark
			}

			watermarkPath := currDir + "\\" + watermarkS
			//сама функция наложения watermark
			funcErr := addWaterMark(paths[i].Path, watermarkPath)
			//так как обработка конкурентна возьмем мьютекс для того чтобы если
			//произошла хоть 1 ошибка то она точно вернется нам
			if funcErr != nil {
				mu.Lock()
				if err == nil {
					err = funcErr
				}
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()
	return err
}

func addWaterMark(bgImg, watermark string) error {
	/*наложение картинки на картину, так же здесь можно указывать
	её прозрачность, угол поворота, пропорцию, положение по горизонтали/вертикали
	*/
	result, err := imagewatermark.ProcessImageWithWatermark(imagewatermark.WatermarkConfig{
		InputPath:             bgImg,
		WatermarkPath:         watermark,
		OpacityAlpha:          0.5,
		WatermarkWidthPercent: 40,
		VerticalAlign:         imagewatermark.VerticalRandom,
		HorizontalAlign:       imagewatermark.HorizontalRandom,
		Spacing:               10,
		RotationDegrees:       20,
	})

	if err != nil {
		return err
	}
	//сохранение на диск
	err = imaging.Save(result, bgImg)
	return nil
}
