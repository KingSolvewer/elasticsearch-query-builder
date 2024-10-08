package main

import (
	"fmt"
	"log"
	"main/ela"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	es := ela.NewYingyanEs()

	//result, err := es.Select(ela.Title, ela.CategoryId, ela.PublishTime, ela.CreateTime).Where(ela.Stat, ela.StatWaitFilter).Size(1).Get()
	result, err := es.Where(ela.Stat, ela.StatWaitFilter).Size(1).Get()
	fmt.Println(result, err)

}
