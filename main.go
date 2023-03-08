package main

import (
	// "time"

	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"rabbit-mq/rabbitmq"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Box Layout")

	task1 := binding.NewString()
	task2 := binding.NewString()

	task1_display := widget.NewLabelWithData(task1)
	task2_display := widget.NewLabelWithData(task2)

	current_task := binding.NewString()
	current_task.Set("Task 1")

	task_content := binding.NewString()

	task1.AddListener(binding.NewDataListener(func() {
		task, err := current_task.Get()
		if task == "Task 1" && err == nil {
			task_content.Set(task1_display.Text)
		}
	}))

	task2.AddListener(binding.NewDataListener(func() {
		task, err := current_task.Get()
		if task == "Task 2" && err == nil {
			task_content.Set(task2_display.Text)
		}
	}))

	textField := widget.NewMultiLineEntry()
	textField.Bind(task_content)

	textField.SetMinRowsVisible(10)

	textField.OnChanged = func(s string) {
		// Send the text to the rabbitmq according to the set up task
		content, err := current_task.Get()
		if err != nil {
			log.Printf("Error: %s", err)
		}
		if content == "Task 1" {
			go rabbitmq.Send("Task 1", s)
		}
		if content == "Task 2" {
			go rabbitmq.Send("Task 2", s)
		}
	}

	go rabbitmq.Recv(&task1, &task2)

	button_1 := widget.NewButton("Task 1", func() {
		current_task.Set("Task 1")
		data, err := task1.Get()
		rabbitmq.FailOnError(err, "Failed to get data for text field")
		err = task_content.Set(data)
		rabbitmq.FailOnError(err, "Failed to set data for text field")
	})

	button_2 := widget.NewButton("Task 2", func() {
		current_task.Set("Task 2")
		data, err := task2.Get()
		rabbitmq.FailOnError(err, "Failed to get data for text field")
		task_content.Set(data)
		err = task_content.Set(data)
		rabbitmq.FailOnError(err, "Failed to set data for text field")
	})

	top_bar := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), button_1, button_2, layout.NewSpacer())
	content := container.New(layout.NewVBoxLayout(), top_bar, layout.NewSpacer(), textField)
	display := container.New(layout.NewVBoxLayout(), canvas.NewText("task 1", color.White), task1_display, canvas.NewText("task 2", color.White), task2_display)

	myWindow.SetTitle("Great")
	myWindow.Resize(fyne.NewSize(400, 400))
	ctx := widget.NewLabelWithData(current_task)
	ctx.TextStyle.Bold = true
	main_container := container.New(layout.NewVBoxLayout(), ctx, content, display)
	myWindow.SetContent(main_container)
	myWindow.ShowAndRun()

	// forever := make(chan int, 1)
	// go src.Send()
	// go src.Recv(forever)
}
