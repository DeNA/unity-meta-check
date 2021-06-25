package inputs

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func StubReadEventPayload(payload *PushOrPullRequestEventPayload, err error) ReadEventPayloadFunc {
	return func(_ typedpath.RawPath) (*PushOrPullRequestEventPayload, error) {
		return payload, err
	}
}
