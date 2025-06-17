package keeper

import (
	"context"
	"encoding/hex"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"strings"
)

type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

func (k Querier) ClusterRelationship(c context.Context, params *types.QueryClusterRelationshipParams) (*types.ClusterRelationshipResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(params.Address)
	if err != nil {
		return nil, err
	}
	accountI := k.accountKeeper.GetAccount(ctx, addr)
	if accountI == nil {
		return nil, core.SignAccountError
	}
	err = verifySignature(params.Signature, params.ClusterId, accountI)
	if err != nil {
		return nil, err
	}
	clusterId := params.ClusterId

	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}
	var allClusterIds []string
	response := &types.ClusterRelationshipResponse{}
	if clusterId == "" {
		personInfo, err := k.GetPersonClusterInfo(ctx, params.Address)
		if err != nil {
			return nil, err
		}
		for key, _ := range personInfo.Owner {
			allClusterIds = append(allClusterIds, key)
			cluster, err := k.GetCluster(ctx, key)
			if err != nil {
				return nil, err
			}
			relationCluster := &types.ClusterRelationship{
				ClusterId:    cluster.ClusterId,
				ClusterName:  cluster.ClusterName,
				ClusterOwner: cluster.ClusterOwner,
			}
			response.ClusterRelationship = append(response.ClusterRelationship, relationCluster)
		}
		if len(allClusterIds) == 0 {
			return nil, nil
		}
	} else {
		clusters := k.GetAllClusters(ctx)
		for _, cluster := range clusters {
			if cluster.ClusterLeader == clusterId {
				relationCluster := &types.ClusterRelationship{
					ClusterId:    cluster.ClusterId,
					ClusterName:  cluster.ClusterName,
					ClusterOwner: cluster.ClusterOwner,
				}
				response.ClusterRelationship = append(response.ClusterRelationship, relationCluster)
			}
		}
	}

	return response, nil
}

func verifySignature(sig string, data string, account authTypes.AccountI) error {
	if sig == "" {
		return core.SignError
	}
	if data == "" {
		data = "root"
	}
	signByte, err := hex.DecodeString(sig)
	if err != nil {
		return err
	}
	ok := account.GetPubKey().VerifySignature([]byte(data), signByte)
	if !ok {
		return core.SignError
	}
	return nil
}

func (k Querier) Params(c context.Context, _ *types.QueryParams) (*types.QueryDaoParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryDaoParamsResponse{Params: params}, nil
}
