package data

import (
	"context"

	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/log"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (p Provider) XReceiptCount() (*hexutil.Big, bool, error) {
	ctx := context.Background()
	query, err := p.EntClient.Receipt.Query().
		Count(ctx)
	if err != nil {
		log.Error(ctx, "Graphql provider err", err)
		return nil, false, err
	}

	big, err := hexutil.DecodeBig(hexutil.EncodeUint64(uint64(query)))
	if err != nil {
		return nil, false, errors.Wrap(err, "decoding block count")
	}

	b := hexutil.Big(*big)

	return &b, true, nil
}