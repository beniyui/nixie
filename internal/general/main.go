package general

import (
	"fmt"
	"nixie/internal/config"
	"os"
)

func general()error{

	//boot
	err := boot()
	if err != nil {
		return err
	}

	//welcomeMSG
	fmt.Printf("welcome to:\n%s\nめんどくさがりだけど優しいボット\n", config.ArtsCont.TitleBanner)
	fmt.Printf("version: %s | exeDir: %s\n", config.EnvCont.Version, config.EnvCont.ExeDir)
	fmt.Printf("now nixie UP. logs:\n=================\n\n")



	//stay alive
	stay()

	//exit
	fmt.Println("bye")
	os.Exit(0)
	return nil
}
