package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	mempl "github.com/tendermint/tendermint/mempool"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	"github.com/tendermint/tendermint/types"
)

//-----------------------------------------------------------------------------
// NOTE: tx should be signed, but this is only checked at the app level (not by Tendermint!)

// BroadcastTxAsync
// More: https://docs.tendermint.com/v0.34/rpc/#/Tx/broadcast_tx_async
func BroadcastTxAsync(ctx *rpctypes.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	err := env.Mempool.CheckTx(tx, nil, mempl.TxInfo{})

	if err != nil {
		return nil, err
	}
	return &ctypes.ResultBroadcastTx{Hash: tx.Hash()}, nil
}

// More: https://docs.tendermint.com/v0.34/rpc/#/Tx/broadcast_tx_sync
func BroadcastTxSync(ctx *rpctypes.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	resCh := make(chan *abci.Response, 1)
	err := env.Mempool.CheckTx(tx, func(res *abci.Response) {
		select {
		case <-ctx.Context().Done():
		case resCh <- res:
		}

	}, mempl.TxInfo{})
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Context().Done():
		return nil, fmt.Errorf("broadcast confirmation not received: %w", ctx.Context().Err())
	case res := <-resCh:
		r := res.GetCheckTx()
		return &ctypes.ResultBroadcastTx{
			Code:      r.Code,
			Data:      r.Data,
			Log:       r.Log,
			Codespace: r.Codespace,
			Hash:      tx.Hash(),
		}, nil
	}
}

// More: https://docs.tendermint.com/v0.34/rpc/#/Tx/broadcast_tx_commit
func BroadcastTxCommit(ctx *rpctypes.Context, tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	subscriber := ctx.RemoteAddr()

	if env.EventBus.NumClients() >= env.Config.MaxSubscriptionClients {
		return nil, fmt.Errorf("max_subscription_clients %d reached", env.Config.MaxSubscriptionClients)
	} else if env.EventBus.NumClientSubscriptions(subscriber) >= env.Config.MaxSubscriptionsPerClient {
		return nil, fmt.Errorf("max_subscriptions_per_client %d reached", env.Config.MaxSubscriptionsPerClient)
	}

	subCtx, cancel := context.WithTimeout(ctx.Context(), SubscribeTimeout)
	defer cancel()
	q := types.EventQueryTxFor(tx)
	deliverTxSub, err := env.EventBus.Subscribe(subCtx, subscriber, q)
	if err != nil {
		err = fmt.Errorf("failed to subscribe to tx: %w", err)
		env.Logger.Error("Error on broadcast_tx_commit", "err", err)
		return nil, err
	}
	defer func() {
		if err := env.EventBus.Unsubscribe(context.Background(), subscriber, q); err != nil {
			env.Logger.Error("Error unsubscribing from eventBus", "err", err)
		}
	}()

	// Broadcast tx and wait for CheckTx result
	checkTxResCh := make(chan *abci.Response, 1)
	err = env.Mempool.CheckTx(tx, func(res *abci.Response) {
		select {
		case <-ctx.Context().Done():
		case checkTxResCh <- res:
		}
	}, mempl.TxInfo{})
	if err != nil {
		env.Logger.Error("Error on broadcastTxCommit", "err", err)
		return nil, fmt.Errorf("error on broadcastTxCommit: %v", err)
	}
	select {
	case <-ctx.Context().Done():
		return nil, fmt.Errorf("broadcast confirmation not received: %w", ctx.Context().Err())
	case checkTxResMsg := <-checkTxResCh:
		checkTxRes := checkTxResMsg.GetCheckTx()
		if checkTxRes.Code != abci.CodeTypeOK {
			return &ctypes.ResultBroadcastTxCommit{
				CheckTx:   *checkTxRes,
				DeliverTx: abci.ResponseDeliverTx{},
				Hash:      tx.Hash(),
			}, nil
		}

		select {
		case msg := <-deliverTxSub.Out():
			deliverTxRes := msg.Data().(types.EventDataTx)
			return &ctypes.ResultBroadcastTxCommit{
				CheckTx:   *checkTxRes,
				DeliverTx: deliverTxRes.Result,
				Hash:      tx.Hash(),
				Height:    deliverTxRes.Height,
			}, nil
		case <-deliverTxSub.Cancelled():
			var reason string
			if deliverTxSub.Err() == nil {
				reason = "Tendermint exited"
			} else {
				reason = deliverTxSub.Err().Error()
			}
			err = fmt.Errorf("deliverTxSub was cancelled (reason: %s)", reason)
			env.Logger.Error("Error on broadcastTxCommit", "err", err)
			return &ctypes.ResultBroadcastTxCommit{
				CheckTx:   *checkTxRes,
				DeliverTx: abci.ResponseDeliverTx{},
				Hash:      tx.Hash(),
			}, err
		case <-time.After(env.Config.TimeoutBroadcastTxCommit):
			err = errors.New("timed out waiting for tx to be included in a block")
			env.Logger.Error("Error on broadcastTxCommit", "err", err)
			return &ctypes.ResultBroadcastTxCommit{
				CheckTx:   *checkTxRes,
				DeliverTx: abci.ResponseDeliverTx{},
				Hash:      tx.Hash(),
			}, err
		}
	}
}

// UnconfirmedTxs gets unconfirmed transactions (maximum ?limit entries)
// including their number.
// More: https://docs.tendermint.com/v0.34/rpc/#/Info/unconfirmed_txs
func UnconfirmedTxs(ctx *rpctypes.Context, limitPtr *int) (*ctypes.ResultUnconfirmedTxs, error) {
	// reuse per_page validator
	limit := validatePerPage(limitPtr)

	txs := env.Mempool.ReapMaxTxs(limit)
	return &ctypes.ResultUnconfirmedTxs{
		Count:      len(txs),
		Total:      env.Mempool.Size(),
		TotalBytes: env.Mempool.SizeBytes(),
		Txs:        txs}, nil
}

// NumUnconfirmedTxs gets number of unconfirmed transactions.
// More: https://docs.tendermint.com/v0.34/rpc/#/Info/num_unconfirmed_txs
func NumUnconfirmedTxs(ctx *rpctypes.Context) (*ctypes.ResultUnconfirmedTxs, error) {
	return &ctypes.ResultUnconfirmedTxs{
		Count:      env.Mempool.Size(),
		Total:      env.Mempool.Size(),
		TotalBytes: env.Mempool.SizeBytes()}, nil
}

func TxExisted(ctx *rpctypes.Context, tx types.Tx) (*ctypes.ResultTxExisted, error) {
	res := &ctypes.ResultTxExisted{}
	existed := env.Mempool.TxExisted(tx)
	if existed {
		res.Existed = 1
		return res, nil
	} else {
		res.Existed = 0
		return res, nil
	}
}

func PoolFirstAndLastTxs(ctx *rpctypes.Context) (*ctypes.PoolFirstAndLastTxs, error) {
	res := &ctypes.PoolFirstAndLastTxs{}
	if env.Mempool.Size() == 0 {
		return res, nil
	}
	txs := env.Mempool.ReapMaxTxs(-1)

	for index, v := range txs {
		if index == 0 {
			res.First = v
		} else if index == env.Mempool.Size()-1 {
			res.Last = v
		}
	}

	return res, nil
}

// CheckTx checks the transaction without executing it. The transaction won't
// be added to the mempool either.
// More: https://docs.tendermint.com/v0.34/rpc/#/Tx/check_tx
func CheckTx(ctx *rpctypes.Context, tx types.Tx) (*ctypes.ResultCheckTx, error) {
	res, err := env.ProxyAppMempool.CheckTxSync(abci.RequestCheckTx{Tx: tx})
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultCheckTx{ResponseCheckTx: *res}, nil
}
