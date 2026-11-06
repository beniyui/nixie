package wizard

import (
	"fmt"
	"nixie/internal/tools/ai"
	"nixie/internal/tools/timer"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

func modelChoices() []string {
	// スライスの要素数が分かっている場合は make で容量を確保すると効率的
	s := make([]string, 0, len(ai.Models))

	for _, m := range ai.Models {
		// append(対象スライス, 追加する値) の形式
		s = append(s, string(m.Explanation))
	}
	return s
}

func ModelIDs() []ai.ModelID {
	s := make([]ai.ModelID, 0, len(ai.Models))

	for _, m := range ai.Models {
		// append(対象スライス, 追加する値) の形式
		s = append(s, ai.ModelID(m.ID))
	}
	return s
}

func Setup() (normal ai.ModelID, think ai.ModelID, name ai.PromptName) {

	fmt.Print("\n")
	fmt.Println("使用するプロンプトは？")
	promptList, err := getPromptList()
	if err != nil {
		timer.FatalExit("failed to get promptList", 5, err)
	}
	i, err := Ask(promptList)
	if err != nil {
		timer.FatalExit("something wrong in Ask Wizard", 5, err)
	}
	pr := ai.PromptName(promptList[i])

	fmt.Print("\n")
	fmt.Print("基本的に使うモデルは？(nixie)\n")
	modelIdList := ModelIDs()
	modelChoices := modelChoices()
	i, err = Ask(modelChoices)
	if err != nil {
		timer.FatalExit("プログラムの不具合", 5, err)
	}
	norm := modelIdList[i]

	fmt.Print("思考に使うモデルは？(nixie#)\n")
	i, err = Ask(modelChoices)
	if err != nil {
		timer.FatalExit("プログラムの不具合", 5, err)
	}
	thi := modelIdList[i]

	return norm, thi, pr
}

func Preload(spiner *spinner.Spinner) error {

	s := spiner

	fmt.Println("プリロードしますか？")
	if preload := AskYesOrNo(); preload {
		startTime := time.Now()
	s.Prefix = fmt.Sprintf("preloading '%s'", idList[i])
					s.Prefix = fmt.Sprintf("preloading '%s'", idList[i])
								s.Start()
											if err := ai.PreloadModel(idList[i]); err != nil {
																s.Stop()
																				return err
																							}
																										s.Stop()
											}
		elapsed := int(time.Since(startTime).Seconds())
		fmt.Printf("完了: 時間: %d秒\n", elapsed)
		if elapsed == 0 {
			fmt.Println("読み込み時間が短い：すでに読み込まれていた可能性")
		}
	}
	return nil
}

func Unload(spiner *spinner.Spinner) error {
	s := spiner
	s.Stop()
	idList := ModelIDs()
	for i := range idList {
		s.Prefix = fmt.Sprintf("preloading '%s'", idList[i])
		s.Start()
		if err := ai.UnloadModel(idList[i]); err != nil {
			s.Stop()
			return err
		}
		s.Stop()
	}
	return nil
}

func getPromptList() ([]string, error) {
	entries, err := os.ReadDir("./config/prompt")
	if err != nil {
		return nil, err
	}

	var names []string
	for _, ent := range entries {
		names = append(names, ent.Name())
	}
	return names, nil
}
