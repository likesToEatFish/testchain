#!/bin/bash

set -e

rm -rf ~/.testchaind/validator4

# node 4
mkdir $HOME/.testchaind/validator4

testchaind init validator4 --chain-id testing-1 --home=$HOME/.testchaind/validator4

# testchaind keys add validator4 --keyring-backend=test --home=$HOME/.testchaind/validator4
echo $(cat /Users/donglieu/script/keys/mnemonic4)| testchaind keys add validator4 --recover --keyring-backend=test --home=$HOME/.testchaind/validator4
# cosmos1qvuhm5m644660nd8377d6l7yz9e9hhm9evmx3x

VALIDATOR4_APP_TOML=$HOME/.testchaind/validator4/config/app.toml

# # validator4
sed -i -E 's|tcp://localhost:1317|tcp://localhost:1313|g' $VALIDATOR4_APP_TOML
sed -i -E 's|localhost:9090|localhost:9082|g' $VALIDATOR4_APP_TOML
sed -i -E 's|localhost:9091|localhost:9083|g' $VALIDATOR4_APP_TOML
sed -i -E 's|tcp://0.0.0.0:10337|tcp://0.0.0.0:10377|g' $VALIDATOR4_APP_TOML
sed -i -E 's|minimum-gas-prices = ""|minimum-gas-prices = "0.0001stake"|g' $VALIDATOR4_APP_TOML


VALIDATOR4_CONFIG=$HOME/.testchaind/validator4/config/config.toml

# # validator4
sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26646|g' $VALIDATOR4_CONFIG
sed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26645|g' $VALIDATOR4_CONFIG
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26644|g' $VALIDATOR4_CONFIG
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR4_CONFIG
sed -i -E 's|prometheus = false|prometheus = true|g' $VALIDATOR4_CONFIG
sed -i -E 's|prometheus_listen_addr = ":26660"|prometheus_listen_addr = ":26600"|g' $VALIDATOR4_CONFIG

cp $HOME/.testchaind/validator1/config/genesis.json $HOME/.testchaind/validator4/config/genesis.json

# copy tendermint node id of validator1 to persistent peers of validator2-3
node1=$(testchaind tendermint show-node-id --home=$HOME/.testchaind/validator1)
node2=$(testchaind tendermint show-node-id --home=$HOME/.testchaind/validator2)
node3=$(testchaind tendermint show-node-id --home=$HOME/.testchaind/validator3)

sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$node1@localhost:26656,$node2@localhost:26653,$node3@localhost:26650\"|g" $HOME/.testchaind/validator4/config/config.toml

# testchaind keys show validator4 -a --keyring-backend=test --home=$HOME/.testchaind/validator4

screen -S gaia4 -t gaia4 -d -m testchaind start --home=$HOME/.testchaind/validator4

# sleep 7

# json_string=$(testchaind tendermint show-validator --home=$HOME/.testchaind/validator4)
# PUB_KEY=$(echo $json_string | jq -r '.key')
# jq --arg newKey "$PUB_KEY" '.pubkey.key = $newKey' /Users/donglieu/script/gaia/validator4.json > temp_val.json && mv temp_val.json /Users/donglieu/script/gaia/validator4.json

# testchaind tx staking create-validator /Users/donglieu/script/gaia/validator4.json --chain-id=testing-1 --from=cosmos1qvuhm5m644660nd8377d6l7yz9e9hhm9evmx3x --gas=500000 --keyring-backend=test --home=$HOME/.testchaind/validator4 -y --fees 500000stake

# sleep 7

# testchaind q staking validators

# testchaind tx staking create-validator \
#   --amount=1000000000000000000000stake \
#   --pubkey=$(testchaind tendermint show-validator --home=$HOME/.testchaind/validator4) \
#   --moniker=MONIKER-YAZ \
#   --chain-id=testing-1 \
#   --commission-rate=0.05 \
#   --commission-max-rate=0.10 \
#   --commission-max-change-rate=0.01 \
#   --min-self-delegation=1 \
#   --from=cosmos1qvuhm5m644660nd8377d6l7yz9e9hhm9evmx3x \
#   --identity="" \
#   --website="" \
#   --details="" \
#   --gas=500000 \
#   --keyring-backend=test \
#   --home=$HOME/.testchaind/validator4 \
#   -y