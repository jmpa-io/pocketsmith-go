<!-- markdownlint-disable MD041 MD010 -->
<p align="center">
  <img src="docs/logo.png"/>
</p>

## `pocketsmith-go`

```diff
+ 📚 A Go abstraction over the Pocketsmith API: https://developers.pocketsmith.com/docs.
```

<a href="LICENSE" target="_blank"><img src="https://img.shields.io/github/license/jmpa-io/pocketsmith-go.svg" alt="GitHub License"></a>
[![CI/CD](https://github.com/jmpa-io/pocketsmith-go/actions/workflows/cicd.yml/badge.svg)](https://github.com/jmpa-io/pocketsmith-go/actions/workflows/cicd.yml)
[![Automerge](https://github.com/jmpa-io/pocketsmith-go/actions/workflows/.github/workflows/dependabot-automerge.yml/badge.svg)](https://github.com/jmpa-io/pocketsmith-go/actions/workflows/.github/workflows/dependabot-automerge.yml)
[![Codecov](https://codecov.io/github/jmpa-io/pocketsmith-go/graph/badge.svg)](https://codecov.io/github/jmpa-io/pocketsmith-go)

## `API Coverage`

The following API endpoints are currently covered by this package:

- [ ] [Get the authorized user](https://developers.pocketsmith.com/reference/get_me-1).
- [ ] [Get user](https://developers.pocketsmith.com/reference/get_users-id-1).
- [ ] [Update user](https://developers.pocketsmith.com/reference/put_users-id-1).
- [ ] [Get institution](https://developers.pocketsmith.com/reference/get_institutions-id-1).
- [ ] [Update institution](https://developers.pocketsmith.com/reference/put_institutions-id-1).
- [ ] [Delete institution](https://developers.pocketsmith.com/reference/delete_institutions-id-1).
- [ ] [List institutions in user](https://developers.pocketsmith.com/reference/get_users-id-institutions-1).
- [ ] [Create institution in user](https://developers.pocketsmith.com/reference/post_users-id-institutions-1).
- [ ] [Get account](https://developers.pocketsmith.com/reference/get_accounts-id-1).
- [ ] [Update account](https://developers.pocketsmith.com/reference/put_accounts-id-1).
- [ ] [Delete account](https://developers.pocketsmith.com/reference/delete_accounts-id-1).
- [ ] [List accounts in user](https://developers.pocketsmith.com/reference/get_users-id-accounts-1).
- [ ] [Update display order of accounts in user](https://developers.pocketsmith.com/reference/put_users-id-accounts-1).
- [ ] [Create an account in user](https://developers.pocketsmith.com/reference/post_users-id-accounts-1).
- [ ] [List accounts in institution](https://developers.pocketsmith.com/reference/get_institutions-id-accounts-1).
- [ ] [Get transaction account](https://developers.pocketsmith.com/reference/get_transaction-accounts-id-1).
- [ ] [Update transaction account](https://developers.pocketsmith.com/reference/put_transaction-accounts-id-1)
- [ ] [List transaction accounts in user](https://developers.pocketsmith.com/reference/get_users-id-transaction-accounts-1).
- [ ] [Get a transaction](https://developers.pocketsmith.com/reference/get_transactions-id-1).
- [ ] [Update a transaction](https://developers.pocketsmith.com/reference/put_transactions-id-1).
- [ ] [Delete a transaction](https://developers.pocketsmith.com/reference/delete_transactions-id).
- [ ] [List transactions in user](https://developers.pocketsmith.com/reference/get_users-id-transactions-1).
- [ ] [List transactions in account](https://developers.pocketsmith.com/reference/get_accounts-id-transactions-1).
- [ ] [List transactions in categories](https://developers.pocketsmith.com/reference/get_categories-id-transactions).
- [ ] [List transactions in transaction account](https://developers.pocketsmith.com/reference/get_transaction-accounts-id-transactions-1).
- [ ] [Create a transaction in transaction account](https://developers.pocketsmith.com/reference/post_transaction-accounts-id-transactions-1).
- [ ] [Get category](https://developers.pocketsmith.com/reference/get_categories-id-1).
- [ ] [Update category](https://developers.pocketsmith.com/reference/put_categories-id-1).
- [ ] [Delete category](https://developers.pocketsmith.com/reference/delete_categories-id-1).
- [ ] [List categories in user](https://developers.pocketsmith.com/reference/get_users-id-categories-1).
- [ ] [Create category in user](https://developers.pocketsmith.com/reference/post_users-id-categories-1).
- [ ] [List category rules in user](https://developers.pocketsmith.com/reference/get_users-id-category-rules-1).
- [ ] [Create category rule in category](https://developers.pocketsmith.com/reference/post_categories-id-category-rules-1).
- [ ] [List budget for user](https://developers.pocketsmith.com/reference/get_users-id-budget-1).
- [ ] [Get budget summary for user](https://developers.pocketsmith.com/reference/get_users-id-budget-summary-1).
- [ ] [Get trend anlysis for user](https://developers.pocketsmith.com/reference/get_users-id-trend-analysis-1).
- [ ] [Delete forcast cache for user](https://developers.pocketsmith.com/reference/delete_users-id-forecast-cache).
- [ ] [Get event](https://developers.pocketsmith.com/reference/get_events-id).
- [ ] [Update event](https://developers.pocketsmith.com/reference/put_events-id).
- [ ] [Delete event](https://developers.pocketsmith.com/reference/delete_events-id).
- [ ] [List events in user](https://developers.pocketsmith.com/reference/get_users-id-events).
- [ ] [List events in scenario](https://developers.pocketsmith.com/reference/get_scenarios-id-events).
- [ ] [Create event in scenario](https://developers.pocketsmith.com/reference/post_scenarios-id-events).
- [ ] [Get attachment](https://developers.pocketsmith.com/reference/get_attachments-id-1).
- [ ] [Update attachment](https://developers.pocketsmith.com/reference/put_attachments-id-1).
- [ ] [Delete attachment](https://developers.pocketsmith.com/reference/delete_attachments-id-1).
- [ ] [List attachments in user](https://developers.pocketsmith.com/reference/get_users-id-attachments-1).
- [ ] [Create attachment in user](https://developers.pocketsmith.com/reference/post_users-id-attachments-1).
- [ ] [List attachments in transaction](https://developers.pocketsmith.com/reference/get_transactions-id-attachments-1).
- [ ] [Assign attachment to transaction](https://developers.pocketsmith.com/reference/post_transactions-id-attachments-1).
- [ ] [Unassign attachment in transaction](https://developers.pocketsmith.com/reference/delete_transactions-transaction-id-attachments-attachment-id-1).
- [ ] [List labels in user](https://developers.pocketsmith.com/reference/get_users-id-labels).
- [ ] [List saved searches in user](https://developers.pocketsmith.com/reference/get_users-id-saved-searches).
- [ ] [List currencies](https://developers.pocketsmith.com/reference/get_currencies).
- [ ] [Get currency](https://developers.pocketsmith.com/reference/get_currencies-id).
- [ ] [List time zones](https://developers.pocketsmith.com/reference/get_time-zones).

## `Usage`

// TODO

## `License`

This work is published under the MIT license.

Please see the [`LICENSE`](./LICENSE) file for details.
