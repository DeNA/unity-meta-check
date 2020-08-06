package meta

func StubGUIDGen(guid *GUID, err error) GUIDGen {
	return func() (*GUID, error) {
		return guid, err
	}
}
