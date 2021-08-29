package lite33

import (
	"github.com/meshplus/pier/model/pb"
)

func (lite *Lite33) verifyHeader(h *pb.BlockHeader) (bool, error) {
	// TODO: blocked by signature mechanism implementation of BitXHub
	return true, nil
}
