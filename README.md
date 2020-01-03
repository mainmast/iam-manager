## IAM MANAGER

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