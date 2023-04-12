package main

import (
	// "time"
	"fmt"

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
	// assign the result of function rabbitmq.randomID() to a variable
	randomID := rabbitmq.RandomID()
	fmt.Println("My queue name is: ", randomID)

	myApp := app.New()
	myWindow := myApp.NewWindow("Box Layout")

	task1 := binding.NewString()
	task2 := binding.NewString()

	task1_status := binding.NewString()
	task2_status := binding.NewString()

	task1_display := widget.NewLabelWithData(task1)
	task2_display := widget.NewLabelWithData(task2)

	task1_status_display := widget.NewLabelWithData(task1_status)
	task2_status_display := widget.NewLabelWithData(task2_status)



	current_task := binding.NewString()
	current_task.Set("Task 1")

	task_content := binding.NewString()

	lock1 := false
	lock2 := false

	task1.AddListener(binding.NewDataListener(func() {
		task, err := current_task.Get()
		if task == "Task 1" && err == nil {
			task_content.Set(task1_display.Text)
			// rabbitmq.Lock("Lock 1", randomID)
			// if task1_status then disable text area
			if(task1_status_display.Text != "(Locked by " + randomID + ")"){
				// log the text
				log.Printf(task1_status_display.Text)
				task_content.Set("Task 1 is locked")
				lock1 = true
			}
		}
	}))

	task2.AddListener(binding.NewDataListener(func() {
		task, err := current_task.Get()
		if task == "Task 2" && err == nil {
			task_content.Set(task2_display.Text)
			// rabbitmq.Lock("Lock 2", randomID)
			// if task2_status then disable text area
			if(task2_status_display.Text != "(Locked by " + randomID + ")"){
				task_content.Set("Task 2 is locked")
				lock2 = true
			}
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
			if(lock1 == true){
				lock1 = false
				return
			}
			go rabbitmq.Send("Task 1", s)
		}
		if content == "Task 2" {
			if(lock2 == true){
				lock2 = false
				return
			}
			go rabbitmq.Send("Task 2", s)
		}
	}

	go rabbitmq.Recv(&task1, &task2, randomID)
	go rabbitmq.LockWatch(&task1_status, &task2_status, randomID)

	button_1 := widget.NewButton("Task 1", func() {
		current_task.Set("Task 1")
		data, err := task1.Get()
		rabbitmq.FailOnError(err, "Failed to get data for text field")
		err = task_content.Set(data)
		rabbitmq.FailOnError(err, "Failed to set data for text field")
		if(task1_status_display.Text != "(Locked by " + randomID + ")"){
			task_content.Set("Task 1 is locked")
			// lock1 = true
		}
	})

	button_2 := widget.NewButton("Task 2", func() {
		current_task.Set("Task 2")
		data, err := task2.Get()
		rabbitmq.FailOnError(err, "Failed to get data for text field")
		task_content.Set(data)
		err = task_content.Set(data)
		rabbitmq.FailOnError(err, "Failed to set data for text field")
		if(task2_status_display.Text != "(Locked by " + randomID + ")"){
			task_content.Set("Task 2 is locked")
			lock1 = true
		}
	})

	task1_lock := widget.NewButton("Lock Task 1", func() {
		rabbitmq.Lock("Lock 1", randomID)
		// get current task
		current, err := current_task.Get()
		if err != nil {
			log.Println(err)
		}
		log.Println(current)
		if(current == "Task 1"){
			task_content.Set(task1_display.Text)
			log.Println(task1_display.Text)
		}
		
	})
	task2_lock := widget.NewButton("Lock Task 2", func() {
		rabbitmq.Lock("Lock 2", randomID)
	})

	top_bar := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), button_1, button_2, layout.NewSpacer())
	content := container.New(layout.NewVBoxLayout(), top_bar, layout.NewSpacer(), textField)
	display := container.New(
		layout.NewVBoxLayout(),
		canvas.NewText("task 1", color.White), task1_display, task1_status_display, task1_lock,
		canvas.NewText("task 2", color.White), task2_display, task2_status_display, task2_lock)

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
