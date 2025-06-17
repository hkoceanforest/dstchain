package group

import (
	"github.com/cosmos/cosmos-sdk/types"
	"time"
)

type ProposalNew struct {
	// id is the unique id of the proposal.
	Id uint64 `json:"id"`
	// group_policy_address is the account address of group policy.
	GroupPolicyAddress string `json:"group_policy_address"`
	// metadata is any arbitrary metadata to attached to the proposal.
	Metadata string `json:"metadata"`
	// proposers are the account addresses of the proposers.
	Proposers []string `json:"proposers"`
	// submit_time is a timestamp specifying when a proposal was submitted.
	SubmitTime time.Time `json:"submit_time"`
	// group_version tracks the version of the group at proposal submission.
	// This field is here for informational purposes only.
	GroupVersion uint64 `json:"group_version"`
	// group_policy_version tracks the version of the group policy at proposal submission.
	// When a decision policy is changed, existing proposals from previous policy
	// versions will become invalid with the `ABORTED` status.
	// This field is here for informational purposes only.
	GroupPolicyVersion uint64 `json:"group_policy_version"`
	// status represents the high level position in the life cycle of the proposal. Initial value is Submitted.
	Status ProposalStatus `json:"status"`
	// final_tally_result contains the sums of all weighted votes for this
	// proposal for each vote option. It is empty at submission, and only
	// populated after tallying, at voting period end or at proposal execution,
	// whichever happens first.
	FinalTallyResult TallyResult `json:"final_tally_result"`
	// voting_period_end is the timestamp before which voting must be done.
	// Unless a successfull MsgExec is called before (to execute a proposal whose
	// tally is successful before the voting period ends), tallying will be done
	// at this point, and the `final_tally_result`and `status` fields will be
	// accordingly updated.
	VotingPeriodEnd time.Time `json:"voting_period_end"`
	// executor_result is the final result of the proposal execution. Initial value is NotRun.
	ExecutorResult ProposalExecutorResult `json:"executor_result"`
	// messages is a list of `sdk.Msg`s that will be executed if the proposal passes.
	Messages []types.Msg `json:"messages"`
}
