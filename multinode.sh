#!/bin/bash
set -xeu

# always returns true so set -e doesn't exit if it is not running.
killall testchaind || true
rm -rf $HOME/.testchaind/

# make four gaia directories
mkdir $HOME/.testchaind
cd $HOME/.testchaind/
mkdir $HOME/.testchaind/validator1
mkdir $HOME/.testchaind/validator2
mkdir $HOME/.testchaind/validator3

# init all three validators
testchaind init --chain-id=testing-1 validator1 --home=$HOME/.testchaind/validator1
testchaind init --chain-id=testing-1 validator2 --home=$HOME/.testchaind/validator2
testchaind init --chain-id=testing-1 validator3 --home=$HOME/.testchaind/validator3

# create keys for all three validators
# cosmos1wa3u4knw74r598quvzydvca42qsmk6jrzmgy07
echo "ozone unfold device pave lemon potato omit insect column wise cover hint narrow large provide kidney episode clay notable milk mention dizzy muffin crazy" | testchaind keys add validator1 --recover --keyring-backend=test --home=$HOME/.testchaind/validator1
# cosmos1w7f3xx7e75p4l7qdym5msqem9rd4dyc4752spg
echo "cactus cable glow spray network group butter recycle weather virtual style bronze treat office energy exit ginger salute slot despair maple bus wheel stool" | testchaind keys add validator2 --recover --keyring-backend=test --home=$HOME/.testchaind/validator2
# cosmos1g9v3zjt6rfkwm4s8sw9wu4jgz9me8pn27f8nyc
echo $(cat /Users/donglieu/script/keys/mnemonic3)| testchaind keys add validator3 --recover --keyring-backend=test --home=$HOME/.testchaind/validator3

# create validator node with tokens to transfer to the three other nodes
testchaind genesis add-genesis-account $(testchaind keys show validator1 -a --keyring-backend=test --home=$HOME/.testchaind/validator1) 10000000000000000000000000000000stake,10000000000000000000000000000000gaia --home=$HOME/.testchaind/validator1 
testchaind genesis add-genesis-account $(testchaind keys show validator2 -a --keyring-backend=test --home=$HOME/.testchaind/validator2) 10000000000000000000000000000000stake,10000000000000000000000000000000gaia --home=$HOME/.testchaind/validator1 
testchaind genesis add-genesis-account $(testchaind keys show validator3 -a --keyring-backend=test --home=$HOME/.testchaind/validator3) 10000000000000000000000000000000stake,10000000000000000000000000000000gaia --home=$HOME/.testchaind/validator1
testchaind genesis add-genesis-account $(testchaind keys show validator1 -a --keyring-backend=test --home=$HOME/.testchaind/validator1) 10000000000000000000000000000000stake,10000000000000000000000000000000gaia --home=$HOME/.testchaind/validator2
testchaind genesis add-genesis-account $(testchaind keys show validator2 -a --keyring-backend=test --home=$HOME/.testchaind/validator2) 10000000000000000000000000000000stake,10000000000000000000000000000000gaia --home=$HOME/.testchaind/validator2 
testchaind genesis add-genesis-account $(testchaind keys show validator3 -a --keyring-backend=test --home=$HOME/.testchaind/validator3) 10000000000000000000000000000000stake,10000000000000000000000000000000gaia --home=$HOME/.testchaind/validator2 
testchaind genesis add-genesis-account $(testchaind keys show validator1 -a --keyring-backend=test --home=$HOME/.testchaind/validator1) 10000000000000000000000000000000stake,10000000000000000000000000000000gaia --home=$HOME/.testchaind/validator3 
testchaind genesis add-genesis-account $(testchaind keys show validator2 -a --keyring-backend=test --home=$HOME/.testchaind/validator2) 10000000000000000000000000000000stake,10000000000000000000000000000000gaia --home=$HOME/.testchaind/validator3 
testchaind genesis add-genesis-account $(testchaind keys show validator3 -a --keyring-backend=test --home=$HOME/.testchaind/validator3) 10000000000000000000000000000000stake,10000000000000000000000000000000gaia --home=$HOME/.testchaind/validator3
testchaind genesis gentx validator1 1000000000000000000000stake --keyring-backend=test --home=$HOME/.testchaind/validator1 --chain-id=testing-1
testchaind genesis gentx validator2 1000000000000000000000stake --keyring-backend=test --home=$HOME/.testchaind/validator2 --chain-id=testing-1
testchaind genesis gentx validator3 1000000000000000000000stake --keyring-backend=test --home=$HOME/.testchaind/validator3 --chain-id=testing-1

cp validator2/config/gentx/*.json $HOME/.testchaind/validator1/config/gentx/
cp validator3/config/gentx/*.json $HOME/.testchaind/validator1/config/gentx/
testchaind genesis collect-gentxs --home=$HOME/.testchaind/validator1 
testchaind genesis collect-gentxs --home=$HOME/.testchaind/validator2
testchaind genesis collect-gentxs --home=$HOME/.testchaind/validator3 

cp validator1/config/genesis.json $HOME/.testchaind/validator2/config/genesis.json
cp validator1/config/genesis.json $HOME/.testchaind/validator3/config/genesis.json


# change app.toml values
VALIDATOR1_APP_TOML=$HOME/.testchaind/validator1/config/app.toml
VALIDATOR2_APP_TOML=$HOME/.testchaind/validator2/config/app.toml
VALIDATOR3_APP_TOML=$HOME/.testchaind/validator3/config/app.toml

# validator1
sed -i -E 's|localhost:9090|localhost:9050|g' $VALIDATOR1_APP_TOML
sed -i -E 's|minimum-gas-prices = ""|minimum-gas-prices = "0.0001stake"|g' $VALIDATOR1_APP_TOML

# validator2
sed -i -E 's|tcp://localhost:1317|tcp://localhost:1316|g' $VALIDATOR2_APP_TOML
# sed -i -E 's|0.0.0.0:9090|0.0.0.0:9088|g' $VALIDATOR2_APP_TOML
sed -i -E 's|localhost:9090|localhost:9088|g' $VALIDATOR2_APP_TOML
# sed -i -E 's|0.0.0.0:9091|0.0.0.0:9089|g' $VALIDATOR2_APP_TOML
sed -i -E 's|localhost:9091|localhost:9089|g' $VALIDATOR2_APP_TOML
sed -i -E 's|minimum-gas-prices = ""|minimum-gas-prices = "0.0001stake"|g' $VALIDATOR2_APP_TOML

# validator3
sed -i -E 's|tcp://localhost:1317|tcp://localhost:1315|g' $VALIDATOR3_APP_TOML
# sed -i -E 's|0.0.0.0:9090|0.0.0.0:9086|g' $VALIDATOR3_APP_TOML
sed -i -E 's|localhost:9090|localhost:9086|g' $VALIDATOR3_APP_TOML
# sed -i -E 's|0.0.0.0:9091|0.0.0.0:9087|g' $VALIDATOR3_APP_TOML
sed -i -E 's|localhost:9091|localhost:9087|g' $VALIDATOR3_APP_TOML
sed -i -E 's|minimum-gas-prices = ""|minimum-gas-prices = "0.0001stake"|g' $VALIDATOR3_APP_TOML

# change config.toml values
VALIDATOR1_CONFIG=$HOME/.testchaind/validator1/config/config.toml
VALIDATOR2_CONFIG=$HOME/.testchaind/validator2/config/config.toml
VALIDATOR3_CONFIG=$HOME/.testchaind/validator3/config/config.toml


# validator1
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR1_CONFIG
# sed -i -E 's|prometheus = false|prometheus = true|g' $VALIDATOR1_CONFIG


# validator2
sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26655|g' $VALIDATOR2_CONFIG
sed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26654|g' $VALIDATOR2_CONFIG
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26653|g' $VALIDATOR2_CONFIG
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR2_CONFIG
sed -i -E 's|prometheus = false|prometheus = true|g' $VALIDATOR2_CONFIG
sed -i -E 's|prometheus_listen_addr = ":26660"|prometheus_listen_addr = ":26630"|g' $VALIDATOR2_CONFIG

# validator3
sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26652|g' $VALIDATOR3_CONFIG
sed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26651|g' $VALIDATOR3_CONFIG
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26650|g' $VALIDATOR3_CONFIG
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR3_CONFIG
sed -i -E 's|prometheus = false|prometheus = true|g' $VALIDATOR3_CONFIG
sed -i -E 's|prometheus_listen_addr = ":26660"|prometheus_listen_addr = ":26620"|g' $VALIDATOR3_CONFIG

# # update

# copy validator1 genesis file to validator2-3
cp $HOME/.testchaind/validator1/config/genesis.json $HOME/.testchaind/validator2/config/genesis.json
cp $HOME/.testchaind/validator1/config/genesis.json $HOME/.testchaind/validator3/config/genesis.json

# copy tendermint node id of validator1 to persistent peers of validator2-3
node1=$(testchaind tendermint show-node-id --home=$HOME/.testchaind/validator1)
node2=$(testchaind tendermint show-node-id --home=$HOME/.testchaind/validator2)
node3=$(testchaind tendermint show-node-id --home=$HOME/.testchaind/validator3)
sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$node1@localhost:26656,$node2@localhost:26653,$node3@localhost:26650\"|g" $HOME/.testchaind/validator1/config/config.toml
sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$node1@localhost:26656,$node2@localhost:26653,$node3@localhost:26650\"|g" $HOME/.testchaind/validator2/config/config.toml
sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$node1@localhost:26656,$node2@localhost:26653,$node3@localhost:26650\"|g" $HOME/.testchaind/validator3/config/config.toml


# # start all three validators/
# testchaind start --home=$HOME/.testchaind/validator1
screen -S gaia1 -t gaia1 -d -m testchaind start --home=$HOME/.testchaind/validator1
screen -S gaia2 -t gaia2 -d -m testchaind start --home=$HOME/.testchaind/validator2
screen -S gaia3 -t gaia3 -d -m testchaind start --home=$HOME/.testchaind/validator3
# testchaind start --home=$HOME/.testchaind/validator3

sleep 7

testchaind q staking validators
testchaind keys list --keyring-backend=test --home=$HOME/.testchaind/validator1
testchaind keys list --keyring-backend=test --home=$HOME/.testchaind/validator2
testchaind keys list --keyring-backend=test --home=$HOME/.testchaind/validator3


# testchaind in-place-testnet testing-1 cosmosvaloper1w7f3xx7e75p4l7qdym5msqem9rd4dyc4mq79dm --home $HOME/.testchaind/validator1
sleep 30
killall testchaind || true
# newOperatorAddress: cosmosvaloper1wa3u4knw74r598quvzydvca42qsmk6jr80u3rd
testchaind in-place-testnet testing-1 cosmosvaloper1wa3u4knw74r598quvzydvca42qsmk6jr80u3rd --home $HOME/.testchaind/validator1  --accounts-to-fund="cosmos1f7twgcq4ypzg7y24wuywy06xmdet8pc4473tnq,cosmos1qvuhm5m644660nd8377d6l7yz9e9hhm9evmx3x" --skip-confirmation
