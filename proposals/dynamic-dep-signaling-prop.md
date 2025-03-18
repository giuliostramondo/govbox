This signaling proposal aims to gather community feedback on updating the `x/gov` module to implement a dynamic proposal deposit mechanism. This feature would replace the current static `MinDeposit` and `MinInitialDeposit` values with an adaptive system that automatically adjusts deposits based on governance activity. If approved, this feature will be included in a future AtomOne software upgrade.

# Motivation

### Addressing Governance Spam

Many Cosmos-based chains suffer from governance spam, where users submit low-cost proposals containing misleading information or scams. While frontend filtering and initial deposit requirements exist, dynamically adjusting `MinInitialDeposit` in response to governance activity can provide a more robust solution.

### Preventing Proposal Overload

An excessive number of active governance proposals can overwhelm stakers, leading to reduced voter participation and governance inefficiencies. The proposed mechanism ensures that governance remains focused by dynamically increasing the deposit requirements when too many proposals are active.

### Reducing Manual Adjustments

Currently, adjusting `MinDeposit` requires governance intervention, making it difficult to respond quickly to changes in proposal volume. A self-regulating deposit mechanism eliminates the need for frequent governance proposals to modify deposit parameters.

# Implementation

The `x/gov` module will be updated to replace the fixed `MinDeposit` and `MinInitialDeposit` parameters with dynamic (independently updated) values determined by the following formula:
 
$$
D_{t+1} = \max(D_{\min}, D_t \times (1+ sign(n_t - N) \times \alpha \times \sqrt[k]{| n_t - N |}))
$$

$$
\alpha = \begin{cases} \alpha_{up} & n_t \gt N \\
\alpha_{down} & n_t \leq N
\end{cases}
$$

$$
sign(n_t - N) = \begin{cases} 1 & n_t \geq N \\
-1 & n_t \lt N
\end{cases}
$$

$$k \in {1, 2, 3, ...}$$
$$0 \lt \alpha_{down} \lt 1$$
$$0 \lt \alpha_{up} \lt 1$$
$$\alpha_{down} \lt \alpha_{up}$$

[See formula in the ADR](https://github.com/atomone-hub/atomone/blob/8561b1e839bf5ac1b27d83b753f89772192502b7/docs/architecture/adr-003-governance-proposal-deposit-auto-throttler.md#decision)

The mechanism updates dynamically:

- When proposals enter or exit the voting/deposit period.
- At regular time intervals (ticks), allowing deposits to gradually adjust even when proposal counts remain stable.

### Key Module Changes

1. **Deprecation of Fixed Deposit Values**
	- `MinDeposit` and `MinInitialDepositRatio` will be deprecated. Attempting to set these parameters in the `x/gov` module will result in an error.
2. **New Dynamic Deposit Parameters:**
The following parameters will be available to fine tune both `MinDeposit` and `MinInitialDeposit`, with each deposit type having their separate collection:
	- `floor_value`: Minimum possible deposit requirement.
	- `update_period`: Time interval for recalculating deposit requirements.
	- `target_active_proposals`: The ideal number of active proposals the system aims to maintain.
	- `increase_ratio` / `decrease_ratio`: Defines how fast deposits adjust.
	- `sensitivity_target_distance`: Controls the steepness of deposit adjustments based on deviations from the target.

3. **Lazily Computed Deposits:**
	- The deposit values are only updated when needed (when proposals change state), and computed lazily for time-based updates rather than continuously updating the value in the blockchain state.

# Testing and Testnet

The dynamic deposit feature has been rigorously tested via unit and end-to-end tests.

The AtomOne public testnet will undergo a coordinated upgrade that mirrors the potential future v2 upgrade of AtomOne and will include the dynamic deposit feature, allowing more testing before mainnet deployment.

# Potential risks

### Increased Complexity

Automating deposit adjustments adds computational complexity. However, lazy evaluation ensures efficient updates without unnecessary overhead.

## User Experience Challenges

Users may find it harder to predict the deposit amount required for a proposal. This can be mitigated with clear client-side tools that display real-time deposit requirements estimates.

# Audit

An audit covering the entire AtomOne codebase and the x/gov module including the dynamic deposit has started in February 2025 and has been recently completed as of today (March 12th), with no findings. More details, including the full audit report, will be released prior to the potential software upgrade proposal that would deploy this functionality on the AtomOne chain.

# Upgrade process

The implementation of the dynamic deposit is contingent upon the successful completion of a third-party audit and thorough validation of its functionality. Once these conditions are met, we anticipate releasing this feature as part of AtomOne v2. The v2 upgrade will be initiated through an upgrade governance proposal, which is tentatively scheduled for submission in Q1 2025, pending results from the ongoing audit for more clarity on final timelines.

Upgrade instructions for validators are documented [here](https://github.com/atomone-hub/atomone/blob/main/UPGRADING.md).

# Codebase

- [https://github.com/atomone-hub/atomone/blob/main/docs/architecture/adr-003-governance-proposal-deposit-auto-throttler.md](https://github.com/atomone-hub/atomone/blob/8561b1e839bf5ac1b27d83b753f89772192502b7/docs/architecture/adr-003-governance-proposal-deposit-auto-throttler.md)
- [https://github.com/atomone-hub/atomone/pull/69](https://github.com/atomone-hub/atomone/pull/69)
- [https://github.com/atomone-hub/atomone/pull/65](https://github.com/atomone-hub/atomone/pull/65)
- [Upgrade instructions](https://github.com/atomone-hub/atomone/blob/main/UPGRADING.md)

# Voting options

- Yes: You are in favor of introducing a dynamic deposit for governance proposals
- No: You are agains having a dynamic proposal deposits
- ABSTAIN - You wish to contribute to the quorum but you formally decline to vote either for or against the proposal.
