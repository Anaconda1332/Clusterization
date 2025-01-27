package myfunc

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"time"

	"parprog/alg"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"math"

	"github.com/fogleman/gg"
)

func ProcessData(numClusters int, filePath string, useMultithreading bool, time time.Duration, flow_num int) string {
	resultText := fmt.Sprintf(
		"Количество кластеров: %d\nПуть к файлу: %s\nМногопоточность: %v\nЗатраченное время: %s\nКоличество потоков: %d",
		numClusters, filePath, useMultithreading, time, flow_num,
	)

	return resultText
}

func CreateResultWindow(myApp fyne.App, numClusters int, filePath string, useMultithreading bool, imagePath string, time time.Duration, flow_num int) {
	resultWindow := myApp.NewWindow("Результат кластеризации")
	resultWindow.Resize(fyne.NewSize(400, 300))

	var resultText string
	if useMultithreading {
		resultText = ProcessData(numClusters, filePath, useMultithreading, time, flow_num)
	} else {
		resultText = ProcessData(numClusters, filePath, useMultithreading, time, 1)
	}

	resultLabel := widget.NewLabel(resultText)

	imageFile, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Ошибка при открытии изображения:", err)
		return
	}
	defer imageFile.Close()

	img, err := png.Decode(imageFile)
	if err != nil {
		fmt.Println("Ошибка при декодировании изображения:", err)
		return
	}

	fyneImage := canvas.NewImageFromImage(img)
	fyneImage.FillMode = canvas.ImageFillOriginal

	resultContent := container.NewVBox(
		resultLabel,
		container.NewHBox(
			fyneImage,
		),
	)

	resultWindow.SetContent(resultContent)

	resultWindow.Show()
}

func hsvToRGB(h, s, v float64) color.RGBA {
	h = math.Mod(h, 360)
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return color.RGBA{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
		A: 255,
	}
}

func CreateImageWithPoints(points []alg.Point_cluster, filename string) error {
	l := len(points)
	var max_x, max_y float64
	max_x = 0
	max_y = 0
	for i := 0; i < l; i++ {
		max_x = max(max_x, points[i].Point[0])
		max_y = max(max_y, points[i].Point[1])
	}
	width, height := 500, 500

	dc := gg.NewContext(width, height)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0, 0, 0)

	dc.DrawLine(50, float64(height-50), float64(width-50), float64(height-50))
	dc.DrawLine(50, 50, 50, float64(height-50))
	dc.Stroke()

	dc.DrawString("X", float64(width-40), float64(height-40))
	dc.DrawString("Y", 40, 40)

	xStart, xEnd := 50, width-50
	yAxis := height - 50
	for x := xStart; x <= xEnd; x += 50 {
		dc.DrawLine(float64(x), float64(yAxis), float64(x), float64(yAxis-10))
		dc.Stroke()

		label := fmt.Sprintf("%d", x-xStart)
		dc.DrawStringAnchored(label, float64(x), float64(yAxis+20), 0.5, 0.5)
	}

	yStart, yEnd := 50, height-50
	xAxis := 50
	for y := yStart; y <= yEnd; y += 50 {
		dc.DrawLine(float64(xAxis), float64(y), float64(xAxis+10), float64(y))
		dc.Stroke()

		label := fmt.Sprintf("%d", yEnd-y)
		dc.DrawStringAnchored(label, float64(xAxis-20), float64(y), 0.5, 0.5)
	}

	for _, pc := range points {
		x := (400 * pc.Point[0] / max_x) + 50
		y := float64(height) - ((400 * pc.Point[1] / max_y) + 50)

		hue := float64(pc.Cluster) * 360.0 / 10.0
		clusterColor := hsvToRGB(hue, 1.0, 1.0)

		dc.SetColor(clusterColor)

		dc.DrawCircle(x, y, 5)
		dc.Fill()
	}

	return dc.SavePNG(filename)
}
