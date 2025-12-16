package files

type FileHub struct {
	Containers map[*FileRoom]bool
	Register   chan *FileRoom
	Unregister chan *FileRoom
}

func NewFileHub() *FileHub {
	return &FileHub{
		Containers: make(map[*FileRoom]bool),
		Register:   make(chan *FileRoom),
		Unregister: make(chan *FileRoom),
	}
}

func (fH *FileHub) Run() {
	for {
		select {
		case container := <-fH.Register:
			fH.Containers[container] = true
		case container := <-fH.Unregister:
			delete(fH.Containers, container)
		}
	}
}
