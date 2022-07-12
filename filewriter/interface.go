package filewriter

type WriterWriteLine interface {
	WriteData(data []string) error
}

type WriterWriteObject interface {
	WriteData(data []byte) error
	GetData() []byte
}
