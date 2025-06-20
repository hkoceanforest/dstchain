// Code generated by protoc-gen-gogo. DO NOT EDIT.

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

const _ = proto.GoGoProtoPackageIsVersion3 

type ClusterHistoricalRewards struct {
	CumulativeRewardRatio github_com_cosmos_cosmos_sdk_types.DecCoins `protobuf:"bytes,1,rep,name=cumulative_reward_ratio,json=cumulativeRewardRatio,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.DecCoins" json:"cumulative_reward_ratio"`
	ReferenceCount        uint32                                      `protobuf:"varint,2,opt,name=reference_count,json=referenceCount,proto3" json:"reference_count,omitempty"`
	HisReward             github_com_cosmos_cosmos_sdk_types.DecCoins `protobuf:"bytes,3,rep,name=his_reward,json=hisReward,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.DecCoins" json:"his_reward"`
	ReceiveCount          int64                                       `protobuf:"varint,4,opt,name=receive_count,json=receiveCount,proto3" json:"receive_count,omitempty"`
}

func (m *ClusterHistoricalRewards) Reset()         { *m = ClusterHistoricalRewards{} }
func (m *ClusterHistoricalRewards) String() string { return proto.CompactTextString(m) }
func (*ClusterHistoricalRewards) ProtoMessage()    {}
func (*ClusterHistoricalRewards) Descriptor() ([]byte, []int) {
	return fileDescriptor_47a2fcf5e349a9ce, []int{0}
}
func (m *ClusterHistoricalRewards) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ClusterHistoricalRewards) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ClusterHistoricalRewards.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ClusterHistoricalRewards) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClusterHistoricalRewards.Merge(m, src)
}
func (m *ClusterHistoricalRewards) XXX_Size() int {
	return m.Size()
}
func (m *ClusterHistoricalRewards) XXX_DiscardUnknown() {
	xxx_messageInfo_ClusterHistoricalRewards.DiscardUnknown(m)
}

var xxx_messageInfo_ClusterHistoricalRewards proto.InternalMessageInfo

func (m *ClusterHistoricalRewards) GetCumulativeRewardRatio() github_com_cosmos_cosmos_sdk_types.DecCoins {
	if m != nil {
		return m.CumulativeRewardRatio
	}
	return nil
}

func (m *ClusterHistoricalRewards) GetReferenceCount() uint32 {
	if m != nil {
		return m.ReferenceCount
	}
	return 0
}

func (m *ClusterHistoricalRewards) GetHisReward() github_com_cosmos_cosmos_sdk_types.DecCoins {
	if m != nil {
		return m.HisReward
	}
	return nil
}

func (m *ClusterHistoricalRewards) GetReceiveCount() int64 {
	if m != nil {
		return m.ReceiveCount
	}
	return 0
}

type ClusterCurrentRewards struct {
	Rewards github_com_cosmos_cosmos_sdk_types.DecCoins `protobuf:"bytes,1,rep,name=rewards,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.DecCoins" json:"rewards"`
	Period  uint64                                      `protobuf:"varint,2,opt,name=period,proto3" json:"period,omitempty"`
}

func (m *ClusterCurrentRewards) Reset()         { *m = ClusterCurrentRewards{} }
func (m *ClusterCurrentRewards) String() string { return proto.CompactTextString(m) }
func (*ClusterCurrentRewards) ProtoMessage()    {}
func (*ClusterCurrentRewards) Descriptor() ([]byte, []int) {
	return fileDescriptor_47a2fcf5e349a9ce, []int{1}
}
func (m *ClusterCurrentRewards) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ClusterCurrentRewards) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ClusterCurrentRewards.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ClusterCurrentRewards) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClusterCurrentRewards.Merge(m, src)
}
func (m *ClusterCurrentRewards) XXX_Size() int {
	return m.Size()
}
func (m *ClusterCurrentRewards) XXX_DiscardUnknown() {
	xxx_messageInfo_ClusterCurrentRewards.DiscardUnknown(m)
}

var xxx_messageInfo_ClusterCurrentRewards proto.InternalMessageInfo

func (m *ClusterCurrentRewards) GetRewards() github_com_cosmos_cosmos_sdk_types.DecCoins {
	if m != nil {
		return m.Rewards
	}
	return nil
}

func (m *ClusterCurrentRewards) GetPeriod() uint64 {
	if m != nil {
		return m.Period
	}
	return 0
}

type ClusterOutstandingRewards struct {
	Rewards github_com_cosmos_cosmos_sdk_types.DecCoins `protobuf:"bytes,1,rep,name=rewards,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.DecCoins" json:"rewards"`
}

func (m *ClusterOutstandingRewards) Reset()         { *m = ClusterOutstandingRewards{} }
func (m *ClusterOutstandingRewards) String() string { return proto.CompactTextString(m) }
func (*ClusterOutstandingRewards) ProtoMessage()    {}
func (*ClusterOutstandingRewards) Descriptor() ([]byte, []int) {
	return fileDescriptor_47a2fcf5e349a9ce, []int{2}
}
func (m *ClusterOutstandingRewards) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ClusterOutstandingRewards) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ClusterOutstandingRewards.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ClusterOutstandingRewards) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClusterOutstandingRewards.Merge(m, src)
}
func (m *ClusterOutstandingRewards) XXX_Size() int {
	return m.Size()
}
func (m *ClusterOutstandingRewards) XXX_DiscardUnknown() {
	xxx_messageInfo_ClusterOutstandingRewards.DiscardUnknown(m)
}

var xxx_messageInfo_ClusterOutstandingRewards proto.InternalMessageInfo

func (m *ClusterOutstandingRewards) GetRewards() github_com_cosmos_cosmos_sdk_types.DecCoins {
	if m != nil {
		return m.Rewards
	}
	return nil
}

type BurnStartingInfo struct {
	PreviousPeriod uint64                                 `protobuf:"varint,1,opt,name=previous_period,json=previousPeriod,proto3" json:"previous_period,omitempty"`
	Stake          github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=stake,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"stake"`
	Height         uint64                                 `protobuf:"varint,3,opt,name=height,proto3" json:"creation_height"`
}

func (m *BurnStartingInfo) Reset()         { *m = BurnStartingInfo{} }
func (m *BurnStartingInfo) String() string { return proto.CompactTextString(m) }
func (*BurnStartingInfo) ProtoMessage()    {}
func (*BurnStartingInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_47a2fcf5e349a9ce, []int{3}
}
func (m *BurnStartingInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BurnStartingInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BurnStartingInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BurnStartingInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BurnStartingInfo.Merge(m, src)
}
func (m *BurnStartingInfo) XXX_Size() int {
	return m.Size()
}
func (m *BurnStartingInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_BurnStartingInfo.DiscardUnknown(m)
}

var xxx_messageInfo_BurnStartingInfo proto.InternalMessageInfo

func (m *BurnStartingInfo) GetPreviousPeriod() uint64 {
	if m != nil {
		return m.PreviousPeriod
	}
	return 0
}

func (m *BurnStartingInfo) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

type RemainderPool struct {
	CommunityPool github_com_cosmos_cosmos_sdk_types.DecCoins `protobuf:"bytes,1,rep,name=community_pool,json=communityPool,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.DecCoins" json:"community_pool"`
}

func (m *RemainderPool) Reset()         { *m = RemainderPool{} }
func (m *RemainderPool) String() string { return proto.CompactTextString(m) }
func (*RemainderPool) ProtoMessage()    {}
func (*RemainderPool) Descriptor() ([]byte, []int) {
	return fileDescriptor_47a2fcf5e349a9ce, []int{4}
}
func (m *RemainderPool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RemainderPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RemainderPool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RemainderPool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemainderPool.Merge(m, src)
}
func (m *RemainderPool) XXX_Size() int {
	return m.Size()
}
func (m *RemainderPool) XXX_DiscardUnknown() {
	xxx_messageInfo_RemainderPool.DiscardUnknown(m)
}

var xxx_messageInfo_RemainderPool proto.InternalMessageInfo

func (m *RemainderPool) GetCommunityPool() github_com_cosmos_cosmos_sdk_types.DecCoins {
	if m != nil {
		return m.CommunityPool
	}
	return nil
}

func init() {
	proto.RegisterType((*ClusterHistoricalRewards)(nil), "freemasonry.dao.v1.ClusterHistoricalRewards")
	proto.RegisterType((*ClusterCurrentRewards)(nil), "freemasonry.dao.v1.ClusterCurrentRewards")
	proto.RegisterType((*ClusterOutstandingRewards)(nil), "freemasonry.dao.v1.ClusterOutstandingRewards")
	proto.RegisterType((*BurnStartingInfo)(nil), "freemasonry.dao.v1.BurnStartingInfo")
	proto.RegisterType((*RemainderPool)(nil), "freemasonry.dao.v1.RemainderPool")
}

func init() { proto.RegisterFile("dao.proto", fileDescriptor_47a2fcf5e349a9ce) }

var fileDescriptor_47a2fcf5e349a9ce = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x94, 0x41, 0x6f, 0xd3, 0x3e,
	0x18, 0xc6, 0xeb, 0x75, 0xff, 0xfe, 0x55, 0x43, 0x3b, 0x14, 0x18, 0x74, 0x13, 0x4a, 0xab, 0x22,
	0xa0, 0xd2, 0xb4, 0x44, 0x65, 0x57, 0x24, 0xa4, 0x76, 0x07, 0x38, 0x31, 0x85, 0x1b, 0x97, 0xc8,
	0x75, 0xde, 0x26, 0x56, 0x13, 0x3b, 0xb2, 0x9d, 0xb2, 0x7e, 0x03, 0x76, 0xe3, 0x03, 0xf0, 0x09,
	0x38, 0x73, 0xe4, 0x03, 0xec, 0x38, 0x71, 0x42, 0x1c, 0x0a, 0x6a, 0x25, 0x0e, 0x7c, 0x0a, 0x94,
	0xd8, 0x1d, 0x1c, 0xb9, 0x54, 0x9c, 0x62, 0x3f, 0xce, 0xfb, 0x3e, 0xbf, 0xf7, 0x49, 0x64, 0xdc,
	0x8c, 0x88, 0xf0, 0x72, 0x29, 0xb4, 0x70, 0x9c, 0xa9, 0x04, 0xc8, 0x88, 0x12, 0x5c, 0x2e, 0xbc,
	0x52, 0x9e, 0x0f, 0x0f, 0x5d, 0x2a, 0x54, 0x26, 0x94, 0x3f, 0x21, 0x0a, 0xfc, 0xf9, 0x70, 0x02,
	0x9a, 0x0c, 0x7d, 0x2a, 0x18, 0x37, 0x35, 0x87, 0x77, 0x62, 0x11, 0x8b, 0x6a, 0xe9, 0x97, 0x2b,
	0xab, 0x1e, 0x98, 0xaa, 0xd0, 0x1c, 0x98, 0x8d, 0x39, 0xea, 0xff, 0xd8, 0xc1, 0x9d, 0x71, 0x5a,
	0x28, 0x0d, 0xf2, 0x39, 0x53, 0x5a, 0x48, 0x46, 0x49, 0x1a, 0xc0, 0x1b, 0x22, 0x23, 0xe5, 0x5c,
	0x20, 0x7c, 0x8f, 0x16, 0x59, 0x91, 0x12, 0xcd, 0xe6, 0x10, 0xca, 0x4a, 0x0e, 0x25, 0xd1, 0x4c,
	0x74, 0x50, 0xaf, 0x3e, 0xb8, 0xf1, 0xe4, 0xbe, 0x67, 0xbb, 0x95, 0x40, 0x9e, 0x05, 0xf2, 0x4e,
	0x81, 0x8e, 0x05, 0xe3, 0xa3, 0x93, 0xcb, 0x65, 0xb7, 0xf6, 0xe1, 0x5b, 0xf7, 0x28, 0x66, 0x3a,
	0x29, 0x26, 0x1e, 0x15, 0x99, 0x75, 0xb7, 0x8f, 0x63, 0x15, 0xcd, 0x7c, 0xbd, 0xc8, 0x41, 0x6d,
	0x6a, 0x54, 0xb0, 0xff, 0xdb, 0xd1, 0x70, 0x04, 0xa5, 0x9f, 0xf3, 0x18, 0xef, 0x49, 0x98, 0x82,
	0x04, 0x4e, 0x21, 0xa4, 0xa2, 0xe0, 0xba, 0xb3, 0xd3, 0x43, 0x83, 0x56, 0xd0, 0xbe, 0x96, 0xc7,
	0xa5, 0xea, 0xe4, 0x18, 0x27, 0x4c, 0x59, 0xd8, 0x4e, 0x7d, 0x5b, 0x98, 0xcd, 0x84, 0x29, 0xc3,
	0xe7, 0x3c, 0xc0, 0x2d, 0x09, 0x14, 0xca, 0x88, 0x0c, 0xd8, 0x6e, 0x0f, 0x0d, 0xea, 0xc1, 0x4d,
	0x2b, 0x56, 0x58, 0xfd, 0xf7, 0x08, 0xef, 0xdb, 0xa0, 0xc7, 0x85, 0x94, 0xc0, 0xf5, 0x26, 0xe5,
	0x19, 0xfe, 0xdf, 0xc0, 0xaa, 0xed, 0x85, 0xba, 0x71, 0x70, 0xee, 0xe2, 0x46, 0x0e, 0x92, 0x89,
	0xa8, 0x4a, 0x6f, 0x37, 0xb0, 0xbb, 0xfe, 0x5b, 0x84, 0x0f, 0x2c, 0xde, 0xcb, 0x42, 0x2b, 0x4d,
	0x78, 0xc4, 0x78, 0xfc, 0x2f, 0x10, 0xfb, 0x9f, 0x10, 0xbe, 0x35, 0x2a, 0x24, 0x7f, 0xa5, 0x89,
	0xd4, 0x8c, 0xc7, 0x2f, 0xf8, 0xb4, 0xfa, 0xfc, 0xb9, 0x84, 0x39, 0x13, 0x85, 0x0a, 0xed, 0x00,
	0xa8, 0x1a, 0xa0, 0xbd, 0x91, 0xcf, 0x2a, 0xd5, 0x09, 0xf0, 0x7f, 0x4a, 0x93, 0x19, 0x54, 0xf3,
	0x35, 0x47, 0x4f, 0x4b, 0x94, 0xaf, 0xcb, 0xee, 0xa3, 0xbf, 0x43, 0xf9, 0xfc, 0xf1, 0x18, 0xdb,
	0xc9, 0x4e, 0x81, 0x06, 0xa6, 0x95, 0x73, 0x84, 0x1b, 0x09, 0xb0, 0x38, 0xd1, 0x9d, 0x7a, 0xe9,
	0x39, 0xba, 0xfd, 0x73, 0xd9, 0xdd, 0xa3, 0x12, 0xca, 0x1f, 0x93, 0x87, 0xe6, 0x28, 0xb0, 0xaf,
	0xf4, 0x2f, 0x10, 0x6e, 0x05, 0x90, 0x11, 0xc6, 0x23, 0x90, 0x67, 0x42, 0xa4, 0xce, 0x39, 0x6e,
	0x53, 0x91, 0x65, 0x05, 0x67, 0x7a, 0x11, 0xe6, 0x42, 0xa4, 0xdb, 0x0b, 0xb1, 0x75, 0x6d, 0x54,
	0x3a, 0x8f, 0x9e, 0x5d, 0xae, 0x5c, 0x74, 0xb5, 0x72, 0xd1, 0xf7, 0x95, 0x8b, 0xde, 0xad, 0xdd,
	0xda, 0xd5, 0xda, 0xad, 0x7d, 0x59, 0xbb, 0xb5, 0xd7, 0x0f, 0xff, 0xbc, 0x5c, 0x28, 0xf5, 0x27,
	0xa9, 0xa0, 0x33, 0x9a, 0x10, 0xc6, 0xfd, 0x73, 0x3f, 0x22, 0xc2, 0x34, 0x9e, 0x34, 0xaa, 0x5b,
	0xe2, 0xe4, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd6, 0x1c, 0xc5, 0xc7, 0x97, 0x04, 0x00, 0x00,
}

func (m *ClusterHistoricalRewards) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClusterHistoricalRewards) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ClusterHistoricalRewards) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ReceiveCount != 0 {
		i = encodeVarintDao(dAtA, i, uint64(m.ReceiveCount))
		i--
		dAtA[i] = 0x20
	}
	if len(m.HisReward) > 0 {
		for iNdEx := len(m.HisReward) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.HisReward[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintDao(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.ReferenceCount != 0 {
		i = encodeVarintDao(dAtA, i, uint64(m.ReferenceCount))
		i--
		dAtA[i] = 0x10
	}
	if len(m.CumulativeRewardRatio) > 0 {
		for iNdEx := len(m.CumulativeRewardRatio) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.CumulativeRewardRatio[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintDao(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *ClusterCurrentRewards) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClusterCurrentRewards) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ClusterCurrentRewards) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Period != 0 {
		i = encodeVarintDao(dAtA, i, uint64(m.Period))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Rewards) > 0 {
		for iNdEx := len(m.Rewards) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Rewards[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintDao(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *ClusterOutstandingRewards) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClusterOutstandingRewards) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ClusterOutstandingRewards) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Rewards) > 0 {
		for iNdEx := len(m.Rewards) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Rewards[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintDao(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *BurnStartingInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BurnStartingInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BurnStartingInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Height != 0 {
		i = encodeVarintDao(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x18
	}
	{
		size := m.Stake.Size()
		i -= size
		if _, err := m.Stake.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDao(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.PreviousPeriod != 0 {
		i = encodeVarintDao(dAtA, i, uint64(m.PreviousPeriod))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *RemainderPool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RemainderPool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RemainderPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.CommunityPool) > 0 {
		for iNdEx := len(m.CommunityPool) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.CommunityPool[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintDao(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintDao(dAtA []byte, offset int, v uint64) int {
	offset -= sovDao(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ClusterHistoricalRewards) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.CumulativeRewardRatio) > 0 {
		for _, e := range m.CumulativeRewardRatio {
			l = e.Size()
			n += 1 + l + sovDao(uint64(l))
		}
	}
	if m.ReferenceCount != 0 {
		n += 1 + sovDao(uint64(m.ReferenceCount))
	}
	if len(m.HisReward) > 0 {
		for _, e := range m.HisReward {
			l = e.Size()
			n += 1 + l + sovDao(uint64(l))
		}
	}
	if m.ReceiveCount != 0 {
		n += 1 + sovDao(uint64(m.ReceiveCount))
	}
	return n
}

func (m *ClusterCurrentRewards) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Rewards) > 0 {
		for _, e := range m.Rewards {
			l = e.Size()
			n += 1 + l + sovDao(uint64(l))
		}
	}
	if m.Period != 0 {
		n += 1 + sovDao(uint64(m.Period))
	}
	return n
}

func (m *ClusterOutstandingRewards) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Rewards) > 0 {
		for _, e := range m.Rewards {
			l = e.Size()
			n += 1 + l + sovDao(uint64(l))
		}
	}
	return n
}

func (m *BurnStartingInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PreviousPeriod != 0 {
		n += 1 + sovDao(uint64(m.PreviousPeriod))
	}
	l = m.Stake.Size()
	n += 1 + l + sovDao(uint64(l))
	if m.Height != 0 {
		n += 1 + sovDao(uint64(m.Height))
	}
	return n
}

func (m *RemainderPool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.CommunityPool) > 0 {
		for _, e := range m.CommunityPool {
			l = e.Size()
			n += 1 + l + sovDao(uint64(l))
		}
	}
	return n
}

func sovDao(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDao(x uint64) (n int) {
	return sovDao(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ClusterHistoricalRewards) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDao
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ClusterHistoricalRewards: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClusterHistoricalRewards: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CumulativeRewardRatio", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDao
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDao
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDao
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CumulativeRewardRatio = append(m.CumulativeRewardRatio, types.DecCoin{})
			if err := m.CumulativeRewardRatio[len(m.CumulativeRewardRatio)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReferenceCount", wireType)
			}
			m.ReferenceCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDao
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ReferenceCount |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HisReward", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDao
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDao
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDao
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HisReward = append(m.HisReward, types.DecCoin{})
			if err := m.HisReward[len(m.HisReward)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReceiveCount", wireType)
			}
			m.ReceiveCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDao
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ReceiveCount |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipDao(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDao
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ClusterCurrentRewards) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDao
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ClusterCurrentRewards: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClusterCurrentRewards: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rewards", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDao
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDao
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDao
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Rewards = append(m.Rewards, types.DecCoin{})
			if err := m.Rewards[len(m.Rewards)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Period", wireType)
			}
			m.Period = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDao
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Period |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipDao(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDao
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ClusterOutstandingRewards) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDao
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ClusterOutstandingRewards: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClusterOutstandingRewards: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rewards", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDao
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDao
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDao
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Rewards = append(m.Rewards, types.DecCoin{})
			if err := m.Rewards[len(m.Rewards)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDao(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDao
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *BurnStartingInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDao
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BurnStartingInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BurnStartingInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PreviousPeriod", wireType)
			}
			m.PreviousPeriod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDao
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PreviousPeriod |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stake", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDao
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDao
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDao
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Stake.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDao
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipDao(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDao
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RemainderPool) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDao
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RemainderPool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RemainderPool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CommunityPool", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDao
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDao
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDao
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CommunityPool = append(m.CommunityPool, types.DecCoin{})
			if err := m.CommunityPool[len(m.CommunityPool)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDao(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDao
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipDao(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDao
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowDao
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowDao
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthDao
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDao
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDao
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDao        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDao          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDao = fmt.Errorf("proto: unexpected end of group")
)
