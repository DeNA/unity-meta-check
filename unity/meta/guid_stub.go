package meta

func AnyGUID() *GUID {
	guid, err := NewGUID([]byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf})
	if err != nil {
		panic(err.Error())
	}
	return guid
}

func ZeroGUID() *GUID {
	guid, err := NewGUID(make([]byte, GUIDByteLength))
	if err != nil {
		panic(err.Error())
	}
	return guid
}
