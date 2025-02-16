package main

import (
	"context"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	flagroll "github.com/novychok/flagroll/gosdk"
)

var (
	apiKey             = "ser.NzRlZWNlM2ItMGExNS00ZTUxLWIyZGQtMDM4YzRmNjcyMjkxEgKehiaulve1BUzvM92bwfq3GFO7IcpYhp6jTszRtgxtwktiYDiEiiRoZntUYEXunLEjIqzx02PDJsnW"
	firstButtonStatus  = false
	secondButtonStatus = false
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ff, err := flagroll.NewClient(apiKey)
	if err != nil {
		fmt.Println(err)
	}

	feature1, err := ff.GetFeature(ctx, "test1")
	if err != nil {
		fmt.Println(err)
	}
	feature2, err := ff.GetFeature(ctx, "test2")
	if err != nil {
		fmt.Println(err)
	}

	featureChan := make(chan flagroll.FeatureFlag)

	go func() {
		err := ff.GetFeatureStatusRealtime(context.Background(), featureChan, feature1, feature2)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		select {
		case f, ok := <-featureChan:
			if !ok {
				fmt.Println("Channel closed. Exiting...")
				return
			}

			if feature1.ID == f.ID && feature1.Name == f.Name {
				firstButtonStatus = f.Active
			} else if feature2.ID == f.ID && feature2.Name == f.Name {
				secondButtonStatus = f.Active
			}

			fmt.Printf("Feature %s is recieved\n", f.ID)
		default:
			// Continue drawing
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

		firstButton(firstButtonStatus)
		secondButton(secondButtonStatus)

		rl.EndDrawing()
	}
}

func firstButton(status bool) {
	color := rl.Blue
	if status {
		color = rl.Yellow
	}

	rl.DrawRectangle(100, 300, 200, 50, color)
	rl.DrawText("Button 1", 150, 315, 20, rl.White)
}

func secondButton(status bool) {
	color := rl.Blue
	if status {
		color = rl.Yellow
	}

	rl.DrawRectangle(400, 300, 200, 50, color)
	rl.DrawText("Button 2", 450, 315, 20, rl.White)
}
