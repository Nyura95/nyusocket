package nyusocket

// Infos public
var Infos = info{}

type info struct {
	List []string
}

func (i *info) add(hash string) {
	i.List = append(i.List, hash)
}

func (i *info) del(hash string) {
	for index, tokenList := range i.List {
		if tokenList == hash {
			i.List = append(i.List[:index], i.List[index+1:]...)
			break
		}
	}
}

func (i *info) Alive(client client) bool {
	for _, tokenList := range i.List {
		if tokenList == client.getHash() {
			return true
		}
	}
	return false
}

func (i *info) NbAlive() int {
	return len(i.List)
}
