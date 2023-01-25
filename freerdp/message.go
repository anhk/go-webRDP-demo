package freerdp

type Bitmap struct {
	X    int    `json:"x"`
	Y    int    `json:"y"`
	W    int    `json:"w"`
	H    int    `json:"h"`
	Data []byte `json:"data"`
}

type Mouse struct {
	Type string `json:"type"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Btn  int    `json:"btn"`
}

type Keyboard struct {
	Type string `json:"type"`
	K    int    `json:"k"`
}

type Message struct {
	Bitmap   *Bitmap   `json:"bitmap,omitempty"`
	Mouse    *Mouse    `json:"mouse"`
	Keyboard *Keyboard `json:"keyboard"`
}
