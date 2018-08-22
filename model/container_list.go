package model

type ContainerList []Container

func (cl ContainerList) Get(i int) Container {
	return cl[i]
}

func (cl ContainerList) FindByHostname(hostname string) int {
	for i, c := range cl {
		if c.Hostname == hostname {
			return i
		}
	}
	return -1
}

func (cl ContainerList) DeleteAt(i int) bool {
	cl = append(cl[:i], cl[i+1:]...)
	return true
}
