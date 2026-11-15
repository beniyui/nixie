package general

import (
	"fmt"
	"nixie/internal/config"
)

func boot() error {
	err := config.Init()
	if err != nil {
		return fmt.Errorf("containers:\n\nprom:\n%+v\nart:\n%+v\nenv\n%+v\n\n===========\n\n%w", config.PromptCont, config.ArtsCont, config.EnvCont, err)
	}

	return nil
}
