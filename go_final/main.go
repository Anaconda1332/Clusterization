package main

import (
	"fmt"
	"os"
	"parprog/alg"
	"parprog/myfunc"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Кластеризация данных")
	myWindow.Resize(fyne.NewSize(500, 300))

	var numClusters, num_flow int
	var filePath string
	var useMultithreading bool

	clusterEntry := widget.NewEntry()
	clusterEntry.SetPlaceHolder("Введите количество кластеров")

	flowEntry := widget.NewEntry()
	flowEntry.SetPlaceHolder("Введите количество потоков")

	fileEntry := widget.NewEntry()
	fileEntry.SetPlaceHolder("Укажите путь к файлу")

	fileDialogButton := widget.NewButton("Выбрать файл", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				filePath = reader.URI().Path()
				fileEntry.SetText(filePath)
			}
		}, myWindow)
	})

	multithreadingRadio := widget.NewRadioGroup([]string{"Да", "Нет"}, func(choice string) {
		useMultithreading = choice == "Да"
	})
	multithreadingRadio.SetSelected("Нет")

	clusterButton := widget.NewButton("Кластеризовать", func() {
		if clusterEntry.Text == "" || fileEntry.Text == "" {
			dialog.ShowInformation("Ошибка", "Заполните все поля", myWindow)
			return
		}

		_, err := fmt.Sscanf(clusterEntry.Text, "%d", &numClusters)
		if err != nil {
			dialog.ShowInformation("Ошибка", "Некорректное количество кластеров", myWindow)
			return
		}

		_, err11 := fmt.Sscanf(flowEntry.Text, "%d", &num_flow)
		if err11 != nil {
			dialog.ShowInformation("Ошибка", "Некорректное количество потоков", myWindow)
			return
		}

		if _, err := os.Stat(fileEntry.Text); os.IsNotExist(err) {
			dialog.ShowInformation("Ошибка", "Файл не существует", myWindow)
			return
		}

		var points []alg.Point_cluster
		var time time.Duration
		if useMultithreading {
			points, time = alg.Go_threaded_clustering(filePath, numClusters, num_flow)
		} else {
			points, time = alg.Go_simple_clustering(filePath, numClusters)
		}

		err1 := myfunc.CreateImageWithPoints(points, "output_points.png")
		if err1 != nil {
			fmt.Println("Ошибка при создании изображения:", err1)
		} else {
			fmt.Println("Изображение успешно создано: output_points.png")
		}

		myfunc.CreateResultWindow(myApp, numClusters, fileEntry.Text, useMultithreading, "D:/go_final/output_points.png", time, num_flow)

	})

	content := container.NewVBox(
		widget.NewLabel("Введите количество кластеров:"),
		clusterEntry,
		widget.NewLabel("Введите количество потоков:"),
		flowEntry,
		widget.NewLabel("Укажите путь к файлу со входными данными:"),
		fileEntry,
		fileDialogButton,
		widget.NewLabel("Хотите ли Вы использовать многопоточность?"),
		multithreadingRadio,
		clusterButton,
	)

	myWindow.SetContent(content)

	myWindow.ShowAndRun()
}
