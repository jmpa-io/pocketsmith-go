
# Default PROJECT, if not given by another Makefile.
ifndef PROJECT
PROJECT=pocketsmith-go
endif

# Targets.
example: binary-go-example
run: example
PHONY += run

---: ## ---

# Includes the common Makefile.
# NOTE: this recursively goes back and finds the `.git` directory and assumes
# this is the root of the project. This could have issues when this assumtion
# is incorrect.
include $(shell while [[ ! -d .git ]]; do cd ..; done; pwd)/Makefile.common.mk

