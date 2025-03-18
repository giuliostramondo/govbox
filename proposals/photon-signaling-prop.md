# Signaling proposal to add the PHOTON token

# Proposal

This signaling proposal aims to gather the community interest in adding the PHOTON token to the AtomOne chain in a future software upgrade. As stated in the Constitution, the PHOTON token is intended to be the only fee token, replacing the ATONE token for paying almost all transaction fees.

The following changes to the AtomOne chain are proposed:

- Add the ability to burn ATONEs in exchange for PHOTONs. This will be the native way to obtain PHOTONs, with the initial PHOTON supply being 0. The conversion rate is calculated based on the total supply of ATONE and minted supply of PHOTON, with PHOTONs total supply being capped at 1 billion. There is no in-protocol mechanism to convert PHOTONs back into ATONE tokens.
- Transaction fees will have to be paid in PHOTON, with the only exception being the newly introduced ATONE to PHOTON burn transaction. Non-exempt transactions that are submitted without using PHOTON as fee token will automatically be rejected.

# Motivation

As written in the AtomOne Constitution, the dual-token model provides a separation of functions, where ATONE serves exclusively as the staking and governance token, while PHOTON functions as the fee token. The introduction of PHOTON as a separate fee token allows more freedom to  ATONE  to work purely as such, dynamically inflating between the two designated bounds of 7% and 20% to disincentivize non-staking targeting a 2/3 bond ratio. Whereas, PHOTON has a fixed max supply as it has no reason to inflate.

The proposed separation enhances clarity around the roles of ATONE and PHOTON. It also avoids Proposals like [848](https://www.mintscan.io/cosmos/proposals/848) being discussed or submitted on the AtomOne chain in the future. We invite the community to look at this [article](https://github.com/atomone-hub/genesis/blob/50882cac6ea4e56b6703d7e3325f35073c75aa6b/STAKING_VS_MONEY.md) for more context on Proposal 848 and the related issues.

PHOTON is also intended to be used for ICS payments (when available).

# Implementation

The details of the implementation are documented in [ADR-002](https://github.com/atomone-hub/atomone/blob/551b5ea4ef12b92c87e6d696373d3ce092a10995/docs/architecture/adr-002-photon-token.md), please refer to this document for a complete understanding of how this is implemented in the AtomOne codebase.

In summary, a new module `x/photon` is added, that brings the following features:

- a new `conversion_rate` query that returns the current conversion rate between ATONE and PHOTON.

```
conversion_rate = (photon_max_supply - photon_supply) / atone_supply
```

- a new `MsgMintPhoton` message that accepts an amount of ATONEs. The input amount (X) is converted to PHOTONs (Y) using the current conversion rate. Then the system burns X ATONEs from the transaction signer's account and mints Y PHOTONs to that same account.
- a new ante decorator is added to the ante handler and validates the transaction fees before execution. All transactions must use PHOTON as the fee token and are otherwise rejected with a specific error. Some messages are exempted and can still use ATONE for the fees, such as `MsgMintPhoton`. This exception list is part of the `x/photon` module parameters and can be modified by a parameter change proposal.

### Gas prices

The gas price mechanism remains unchanged - validators still set their own gas prices.

# Testing and Testnet

The `x/photon` module has been rigorously tested via unit and end-to-end tests.

The AtomOne public testnet will undergo a coordinated upgrade that mirrors the potential future v2 upgrade of AtomOne and will include the `x/photon` module, allowing more testing, before the final mainnet deployment. 

# Potential Risks

### User Experience Challenges

Having two tokens complicates the user experience, as users must manage two assets with different roles. While most clients suggest a fee token list to the user when a transaction is submitted, users may be unfamiliar with this process, and in addition selecting the correct fee token might not be intuitive.

This can be mitigated through client updates that automatically suggest the appropriate token for each transaction. As a simple rule, clients should default to PHOTON as the fee token for all transactions except for the `MsgMintPhoton` message, for which clients should default to ATONE.

# Security Audits

A third-party securityaudit covering the existing  AtomOne codebase and the new `x/photon` module was conducted between February and March of 2025. No security issues were identified during the audit.. More details about the audit and the full report will be released prior to the software upgrade proposal that would deploy this functionality on the AtomOne chain.

# Upgrade Process

The implementation of the `x/photon` module is contingent upon the successful completion of a third-party audit and thorough validation of its functionality through the public testnet. Once these conditions are met, we anticipate releasing this feature as part of AtomOne v2. The v2 upgrade will be initiated through an upgrade governance proposal, which is tentatively scheduled for submission in Q1 2025, pending results from the ongoing audit for more clarity on final timelines.

Upgrade instructions for validators are documented [here](https://github.com/atomone-hub/atomone/blob/main/UPGRADING.md)

# Codebase

- [adr-002-photon-token.md](https://github.com/atomone-hub/atomone/blob/551b5ea4ef12b92c87e6d696373d3ce092a10995/docs/architecture/adr-002-photon-token.md)
- https://github.com/atomone-hub/atomone/pull/34
- https://github.com/atomone-hub/atomone/pull/57
- [Upgrade instructions](https://github.com/atomone-hub/atomone/blob/main/UPGRADING.md)

# Voting Options

- Yes: You are in favor of introducing PHOTON as the only fee token.
- No: You are against this dual-token model.
- ABSTAIN - You wish to contribute to the quorum but you formally decline to vote either for or against the proposal.
