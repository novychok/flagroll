package main

import (
	"context"
	"fmt"

	flagroll "github.com/novychok/flagroll/gosdk"
)

var apiKey = "ser.NzRlZWNlM2ItMGExNS00ZTUxLWIyZGQtMDM4YzRmNjcyMjkxEgKehiaulve1BUzvM92bwfq3GFO7IcpYhp6jTszRtgxtwktiYDiEiiRoZntUYEXunLEjIqzx02PDJsnW"

func main() {
	ff, err := flagroll.NewClient(apiKey)
	if err != nil {
		fmt.Println(err)
	}

	isVersionv2, err := ff.GetFeatureStatus(context.Background(), "test123")
	if err != nil {
		fmt.Println(err)
	}

	// showData, err := ff.GetFeatureValue("bobol")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	fmt.Println(isVersionv2)
	// fmt.Println(showData)

}
