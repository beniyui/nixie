package general

import (
	"fmt"
	"nixie/internal/config"
)

func general() error {

	//init; loading the env file and save the containers
	//envCont, artCont......
	err := config.Init()
	if err != nil {
		return fmt.Errorf("config.Init: %v",err)
	}

	//welcomeMSG
	fmt.Printf("welcome to:\n%s\nめんどくさがりだけど優しいボット\n", config.ArtsCont.TitleBanner)
	fmt.Printf("version: %s | exeDir: %s\n", config.EnvCont.Version, config.EnvCont.ExeDir)
	fmt.Printf("now nixie UP. logs:\n=================\n\n")

	//stay alive
	stay()

	//exit
	fmt.Println("bye")
	return nil
}
