package nyusocket

// Infos public
var Infos = info{}

type info struct {
	List []string
}

func (i *info) Add(token string) {
	i.List = append(i.List, token)
}

func (i *info) Del(token string) {
	for index, tokenList := range i.List {
		if tokenList == token {
			i.List = append(i.List[:index], i.List[index+1:]...)
			break
		}
	}
}

func (i *info) Alive(token string) bool {
	for _, tokenList := range i.List {
		if tokenList == token {
			return true
		}
	}
	return false
}

func (i *info) NbAlive() int {
	return len(i.List)
}
