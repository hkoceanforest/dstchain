<!--
order: 2
-->

# Cosmos Coin Registration

Learn how to register native Cosmos Coins through [Evmos Governance](../../users/governance/overview.md). {synopsis}

:::tip
**Note**: Not sure what the difference between Cosmos Coin and ERC-20 Registration is? You're in the right place if you want to add functionality to convert a native Cosmos Coin to an ERC-20 token representation. If an ERC-20 contract corresponding to your token already exists, and you want to add functionality to convert the ERC-20 token to a native Cosmos Coin denomination, check out [ERC-20 Registration](./erc20_registration.md) instead.

Still confused? Learn more about the differences [here](../../../x/erc20/spec/01_concepts.md).
:::

The ERC-20 Module (also known as `x/erc20`) allows users to instantly convert [ERC-20](https://ethereum.org/en/developers/docs/standards/tokens/erc-20) tokens into native Cosmos Coins, and vice versa. This allows users to exchange assets interchangeably in two entirely different layers, the EVM and Cosmos.

Application-wise, the ERC-20 module allows DeFi protocols to seamlessly integrate with Evmos and the Cosmos ecosystem. Using the module, developers can build smart contracts on Evmos and use the generated ERC-20 tokens for other [applications on the Cosmos ecosystem](https://mapofzones.com), such as:

- earning $OSMO staking rewards
- taking part in governance proposals by voting with $ATOM

Registering a native Cosmos Coin means registering a new mapping between a Cosmos Coin denomination and a new ERC-20 token contract, also known as a Token Pair. Token Pairs enable users to convert ERC-20 tokens into their native Cosmos Coin representation (and vice-versa), and can only be created via a governance proposal.

More information can be found in [this blog post](https://medium.com/evmos/introducing-evmos-erc20-module-f40a61e05273), which introduced the ERC-20 Module on Evmos.

To register a Cosmos Coin, consider the following stages:

## Drafting the Cosmos Coin Proposal

The following topics must be addressed when drafting a Cosmos Coin Proposal:

1. Provide the profile of the person(s)/entity making the proposal.

	Who are you? What is your involvement in Cosmos and/or other blockchain networks? If you are working with a team, who are the team members involved and what is their relevant experience? What is the mission statement of your organization or business? Do you have a website? Showcase some work you've done and some proof of who you are.

2. Promote understanding of the ERC-20 Module.

  Make sure to mention the original [blog post](https://medium.com/evmos/introducing-evmos-erc20-module-f40a61e05273) that introduced the ERC-20 Module, along with a brief explanation of what the ERC-20 Module does. It's also a good idea to link the [ERC-20 Module documentation](https://docs.evmos.org/modules/erc20/)!

3. Describe how ERC-20 Module changes will be made.

  Give a breakdown of the proposal's payload, and explain in layman terms what the proposal will do if it passes. Detail precautions taken during contract and proposal formulation, if applicable (including consultations made prior to proposal creation, how contracts were tested, and any third-party reviews). Finally, mention the risks involved in the proposal, depending on the direction of IBC Coin and ERC-20.

Remember to provide links to the relevant [Commonwealth Evmos community](https://commonwealth.im/evmos) discussions concerning your proposal, as well as the [proposal on testnet](#submit-the-proposal-to-the-testnet).

## Adding Network to Evmos Chain Registry

All chain registry information can be found [here](https://github.com/evmos/chain-token-registry). The [Cosmos Chain Registry](https://github.com/cosmos/chain-registry) is used to pull in the list of RPC, gRPC, and REST endpoints to power the Evmos Dashboard. To add chain registry information, please submit a pull request with the network details. Pull requests should be merged one business day after governance passes.

## Submitting the Cosmos Coin Proposal

After the drafting process, the Cosmos Coin Proposal can be submitted.

### Formatting the Proposal's Text

The ideal format for a proposal is as a Markdown file (ie. `.md`) in a Github repo or [HackMd](https://hackmd.io/). Markdown
is a simple and accessible format for writing plain text files that is easy to

<!-- markdown-link-check-disable-next-line -->

learn. See the [Github Markdown Guide](https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax) for details on
writing markdown files.

### Submit the Proposal to Testnet

To [submit the proposal](../../users/governance/submitting.md) to testnet through the command line with [`evmosd`](../../validators/quickstart/binary.md), use the following command with `register-coin`:

```bash
evmosd tx gov submit-proposal register-coin <path/to/metadata.json> \
  --title=<title> \
  --description=<description> \
  --deposit="1000000aevmos" \
  --from=<mykey> \
  --chain-id=<testnet_chain_id> \
  --node <address>
```

where `metadata.json` contains (as an example, for Osmosis):

```json
{
    "metadata": [
        {
			"description": "The native staking and governance token of the Osmosis chain",
			"denom_units": [
				{
						"denom": "ibc/<HASH>",
						"exponent": 0,
						"aliases": ["ibcuosmo"]
				},
				{
						"denom": "OSMO",
						"exponent": 6
				}
			],
			"base": "ibc/<HASH>",
			"display": "OSMO",
			"name": "Osmo",
			"symbol": "OSMO"
		}
	]
}
```

Note that you can register multiple coins in one proposal by extending the `"metadata"` array with multiple metadata entries, e.g.

```json
{
    "metadata": [
        {
			"description": "The native staking and governance token of the Osmosis chain",
            // ...
   			"symbol": "OSMO"
		},
        {
			"description": "The native staking and governance token of the Cosmos chain",
            // ...
   			"symbol": "ATOM"
		},
	]
}
```

However, note that if the CLI is used to create a proposal, and `description` is set using a flag, the text will be [escaped](https://en.wikipedia.org/wiki/Escape_sequences_in_C) which may have undesired effects. If the proposal creator is using markdown or line breaks it's recommended to put the proposal text into a json file and include that file as part of the CLI proposal, as opposed to individual fields in flags. The process of creating a json file containing the proposal can be found [here](../../users/governance/submitting.md#formatting-the-json-file-for-the-governance-proposal), and the CLI command for submitting the file is below:

```bash
evmosd tx gov submit-proposal register-coin --proposal=<path/to/proposal.json>
```

You may want to submit your proposal to the testnet chain before the mainnet for a number of reasons, such as wanting to see what the proposal description will look like, to share what the proposal will look like in advance with stakeholders, and to signal that your proposal is about to go live on the mainnet.

Submitting your proposal to the testnet increases the likelihood of engagement and the possibility that you will be alerted to a flaw before deploying your proposal to mainnet.

## The On-Chain ERC-20 Proposal

A majority of the voting community should probably be aware of the proposal and have considered it before the proposal goes live on-chain. If you're taking a conservative approach, you should have reasonable confidence that your proposal will pass before risking deposit contributions by [submitting the proposal](../../users/governance/submitting.md). Make revisions to your draft proposal after each stage of engagement.

### The Deposit Period

The deposit period currently lasts 14 days. If you submitted your transaction with the minimum deposit (64 EVMOS), your proposal will immediately enter the voting period. If you didn't submit the minimum deposit amount (currently 64 EVMOS), then this may be an opportunity for others to show their support by contributing (and risking) their EVMOS as a bond for your proposal. You can request contributions openly and also contact stakeholders directly (particularly stakeholders who are enthusiastic about your proposal). Remember that each contributor is risking their funds, and you can [read more about the conditions for burning deposits here](../../users/governance/process.md#burned-deposits).

This is a stage where proposals may begin to get broader attention. Most popular explorers currently display proposals that are in the deposit period, but due to proposal spamming, this may change.

A large cross-section of the blockchain/cryptocurrency community exists on Twitter. Having your proposal in the deposit period is a good time to engage the Evmos community to prepare validators to vote and EVMOS-holders that are staking.

### The Voting Period

At this point you'll want to track which validator has voted and which has not. You'll want to re-engage directly with top stake-holders, ie. the highest-ranking validator operators, to ensure that:

1. they are aware of your proposal;
2. they can ask you any questions about your proposal; and
3. they are prepared to vote.

Remember that any voter may change their vote at any time before the voting period ends. That historically doesn't happen often, but there may be an opportunity to convince a voter to change their vote. The biggest risk is that stakeholders won't vote at all (for a number of reasons). Validator operators tend to need multiple reminders to vote. How you choose to contact validator operators, how often, and what you say is up to you--remember that no validator is obligated to vote, and that operators are likely occupied by competing demands for their attention. Take care not to stress any potential relationship with validator operators.
