


//Modelsのexplanation(説明)のスライスを取得する関数 外部向け
func ModelChoices() []string {
	//make で容量を確保すると効率的
	s := make([]string, 0, len(Models))

	for _, m := range Models {
		s = append(s, string(m.Explanation))
	}
	return s
}

//各モデルのIDのスライスを取得する関数 外部向け
func ModelIDs() []string {
	s := make([]string, 0, len(Models))

	for _, m := range Models {
		// append(対象スライス, 追加する値) の形式
		s = append(s, string(m.ID))
	}
	return s
}
