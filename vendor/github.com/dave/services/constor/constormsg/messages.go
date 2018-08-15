package constormsg

import "encoding/gob"

func RegisterTypes() {
	gob.Register(Storing{})
}

type Storing struct {
	Starting  bool
	Finished  int
	Unchanged int
	Remain    int
	Done      bool
}
