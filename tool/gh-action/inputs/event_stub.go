package inputs

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func StubReadEventPayload(payload *PullRequestEventPayload, err error) ReadEventPayloadFunc {
	return func(_ typedpath.RawPath) (*PullRequestEventPayload, error) {
		return payload, err
	}
}
