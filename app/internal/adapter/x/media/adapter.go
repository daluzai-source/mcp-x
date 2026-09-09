package media

import "github.com/michimani/gotwi"

// chunkBytes stays under the 5 MB APPEND limit with room for the multipart
// envelope the library wraps each segment in.
const chunkBytes = 4 << 20

type Adapter struct {
	client gotwi.IClient
}

func New(client gotwi.IClient) *Adapter {
	return &Adapter{client: client}
}

func (a *Adapter) ChunkBytes() int {
	return chunkBytes
}
