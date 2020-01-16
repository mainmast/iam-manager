## Identity & Access Management

This is a very simplistic IAM manager for supporting a platform that plans to allow their customers to have their own segregated data, and manage their own users/accounts.

We use the following terminology assuming your whole product uses this package as it's login/authentication:

- Organisation: This is an organisation that uses your SaaS platform (think your AWS account is an organisation). We support organisation hierarchies.
- IAM User: This is a person who can login to an organisation account (think your AWS login)
- Platform User: This is a person who can login to an organisations SaaS platform (they have no IAM access to YOUR SaaS platform)
- Account: This is a person who can log into an organisations end product, but not their platform.

Controls all resources related with IAM main actors such as 

- Create Organisations
- Create Accounts
- Associate users to Organisations
- Assign roles to users

### Migrations

there are two ways to run migrations `upgrade` and `downgrade` , in order to create a new migration script a file should be created with the next format in `conf` folder:  `version`+`script_name`+`action`+`.sql`  action = `up` to upgrade and = `down` to downgrade.

example

```go

// upgrade migrations
go run ./cmd/iam-migrations --action=upgrade

// downgrade migrations
go run ./cmd/iam-migrations --action=downgrade

```