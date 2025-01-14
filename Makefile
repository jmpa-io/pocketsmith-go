
# Default PROJECT, if not given by another Makefile.
ifndef PROJECT
PROJECT=pocketsmith-go
endif

# Variables.
TOKEN = $(shell aws ssm get-parameter --name "/tokens/pocketsmith" --query 'Parameter.Value' --output text --with-decryption)

# Targets.
accounts: binary-go-accounts ## Builds the 'accounts' binary.
authed-user: binary-go-authed-user ## Builds the 'authed-user' binary.
tracing: binary-go-tracing ## Builds the `tracing` binary.
run: accounts authed-user tracing

PHONY += accounts authed-user tracing run

get-token: ## Retrieves the Pocketsmith token from AWS SSM Parameter Store.
get-token:
	@echo "export POCKETSMITH_TOKEN=\"$(TOKEN)\""

---: ## ---

# Includes the common Makefile.
# NOTE: this recursively goes back and finds the `.git` directory and assumes
# this is the root of the project. This could have issues when this assumtion
# is incorrect.
include $(shell while [[ ! -d .git ]]; do cd ..; done; pwd)/Makefile.common.mk

