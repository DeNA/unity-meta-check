package meta

import "crypto/rand"

type GUIDGen func() (*GUID, error)

func RandomGUIDGenerator() GUIDGen {
	return func() (*GUID, error) {
		bytes := make([]byte, GUIDByteLength)
		_, err := rand.Read(bytes)
		if err != nil {
			return nil, err
		}
		return NewGUID(bytes)
	}
}
