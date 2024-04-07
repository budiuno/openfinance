this services is consist three endpoint

1. "/v1/account/{bank_code}/{account_number}" to check bank account exist or not
2. "/v1/disbursements" to create disbursement
3. "/v1/callback/disbursements" to received callback for disbursement status update

this services use Postgres, and can spawn by just run the docker compose up
for the migration table please run ./run_migrations.sh

"/v1/account/{bank_code}/{account_number}" & "/v1/disbursements" is protect by authentication
token for authentication will show on console when service running, use the token as --header 'Authorization: ${token}'

"/v1/account/{bank_code}/{account_number}" hit mockapi https://66100a360640280f219c2844.mockapi.io/api/v1/accounts as mock of bank api to check account

"/v1/disbursements" hit mockapi https://66100a360640280f219c2844.mockapi.io/api/v1/disburse as mock of bank api to post disbursement

Postman collection is also attached
