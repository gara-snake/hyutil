package hcollection

// Identifiable 識別子取得型
type Identifiable interface {
	ID() string
}

// StrContains 配列の中に要素が存在するかどうか確認する
func StrContains(list []string, that string) bool {
	for _, item := range list {
		if item == that {
			return true
		}
	}
	return false
}

// Contains 配列の中に要素が存在するかどうか確認する
func Contains(list []Identifiable, target Identifiable) bool {
	for _, item := range list {
		if item.ID() == target.ID() {
			return true
		}
	}
	return false
}
