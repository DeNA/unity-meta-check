package inputs

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func StubReadEventPayload(payload *EventPayload, err error) ReadEventPayloadFunc {
	return func(_ typedpath.RawPath) (*EventPayload, error) {
		return payload, err
	}
}
