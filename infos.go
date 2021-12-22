package nyusocket

type Info struct {
	List []string
}

func (i *Info) add(hash string) {
	i.List = append(i.List, hash)
}

func (i *Info) del(hash string) {
	for index, tokenList := range i.List {
		if tokenList == hash {
			i.List = append(i.List[:index], i.List[index+1:]...)
			break
		}
	}
}

func (i *Info) Alive(client client) bool {
	for _, tokenList := range i.List {
		if tokenList == client.getHash() {
			return true
		}
	}
	return false
}

func (i *Info) NbAlive() int {
	return len(i.List)
}
