MONIKER=adminmoduletest1
CHAIN_ID=adminmodule
CHAIN_HOME=$(HOME)/.admin-module
STARPORT_HOME=$(HOME)/.starport
STARPORT_VERSION=0.17.3
COSMOS_SDK_VERSION=$(shell cat go.mod | grep cosmos-sdk | cut -d ' ' -f2 | sed 's/^v//')

ldflags = -X github.com/cosmos/cosmos-sdk/version.AppName=admin-moduled \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(COSMOS_SDK_VERSION)


containerProtoVer=0.9.0
containerProtoImage=ghcr.io/cosmos/proto-builder:$(containerProtoVer)
containerProtoGen=cosmos-sdk-proto-gen-$(containerProtoVer)
containerProtoGenSwagger=cosmos-sdk-proto-gen-swagger-$(containerProtoVer)
containerProtoFmt=cosmos-sdk-proto-fmt-$(containerProtoVer)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGen}$$"; then docker start -a $(containerProtoGen); else docker run --name $(containerProtoGen) -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protocgen.sh; fi

.PHONY: version
version:
	echo $(COSMOS_SDK_VERSION)

.PHONY: .get-starportg
.get-starport:
	[ -f bin/starport ] || $$(curl -LO https://github.com/tendermint/starport/releases/download/v$(STARPORT_VERSION)/starport_$(STARPORT_VERSION)_linux_amd64.tar.gz && tar -xzf starport_$(STARPORT_VERSION)_linux_amd64.tar.gz starport && mv starport bin/ && chmod +x bin/starport && rm starport_$(STARPORT_VERSION)_linux_amd64.tar.gz)

.PHONY: build
build:
	go build -ldflags '$(ldflags)' -o bin/admin-moduled cmd/admin-moduled/main.go

.PHONY: maybebuild
maybebuild:
	[ -f bin/admin-moduled ] || make build

.PHONY: generate
generate: .get-starport
	bin/starport generate proto-go && bin/starport generate openapi

.PHONY: test
test:
	go test ./...

.PHONY: local-clean
local-clean:
	rm -rf $(CHAIN_HOME) && rm -rf $(STARPORT_HOME)

.PHONY: local-keys
local-keys: maybebuild
	bin/admin-moduled keys show alice --home $(CHAIN_HOME) > /dev/null 2>&1 || (sleep 1; echo "earn noble employ useful space craft staff blast exact pluck siren physical biology short suit oval open legend humble pill series devote wealth hungry") | bin/admin-moduled keys add alice --recover --home $(CHAIN_HOME)
	bin/admin-moduled keys show bob --home $(CHAIN_HOME) > /dev/null 2>&1 || (sleep 1; echo "lawn pigeon use festival elder trust wish rose law family about web fiber jealous daughter vote history grant quarter fetch soft poem aware truly") | bin/admin-moduled keys add bob --recover --home $(CHAIN_HOME)
	bin/admin-moduled keys show user1 --home $(CHAIN_HOME) > /dev/null 2>&1 || (sleep 1; echo "hello turn increase august raw toss hurdle craft baby arrow aware shield maple net six math chase debris chase wet benefit rent segment beauty") | bin/admin-moduled keys add user1 --recover --home $(CHAIN_HOME)
	bin/admin-moduled keys show user2 --home $(CHAIN_HOME) > /dev/null 2>&1 || (sleep 1; echo "high return silly coyote skin trumpet stock bicycle enjoy common exact sure") | bin/admin-moduled keys add user2 --recover --home $(CHAIN_HOME)

.PHONY: local-init
local-init: local-keys
	bin/admin-moduled init $(MONIKER) --chain-id $(CHAIN_ID) --home $(CHAIN_HOME)
	bin/admin-moduled add-genesis-account alice 10000000000000000000000001stake --home $(CHAIN_HOME)
	bin/admin-moduled gentx alice 1000000000stake --chain-id $(CHAIN_ID) --home $(CHAIN_HOME)
	bin/admin-moduled collect-gentxs --home $(CHAIN_HOME)
	sed -i "s/prometheus = false/prometheus = true/" $(CHAIN_HOME)/config/config.toml
	sed -i "s/cors_allowed_origins = \[\]/cors_allowed_origins = [\"*\"]/" $(CHAIN_HOME)/config/config.toml
	cat $(CHAIN_HOME)/config/app.toml | tr '\n' '\r' | sed "s/# Enable defines if the API server should be enabled.\renable = false/# Enable defines if the API server should be enabled.\renable = true/" | tr '\r' '\n' > /tmp/app.toml.tmp && mv /tmp/app.toml.tmp $(CHAIN_HOME)/config/app.toml
	sed -i "s/swagger = false/swagger = true/" $(CHAIN_HOME)/config/app.toml
	sed -i "s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/" $(CHAIN_HOME)/config/app.toml
	jq '.app_state.adminmodule.admins[0] = .app_state.auth.accounts[0].address' $(CHAIN_HOME)/config/genesis.json > /tmp/genesis.json.tmp && mv /tmp/genesis.json.tmp $(CHAIN_HOME)/config/genesis.json

.PHONY: local-start
local-start:
	bin/admin-moduled start --home $(CHAIN_HOME)

.PHONY: starport-start
starport-start: .get-starport
	bin/starport chain serve
