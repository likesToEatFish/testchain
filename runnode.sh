#!/bin/bash
killall testchaind || true
rm -rf $HOME/.testchain/

# echo $(cat ./keys/mnemonic1)| testchaind keys add val1 --keyring-backend test --recover
# echo $(cat ./keys/mnemonic2)| testchaind keys add val3 --keyring-backend test  --recover --home /Users/donglieu/script/xion/node4
echo $(cat /Users/donglieu/script/keys/mnemonic1)| testchaind keys add val --keyring-backend test  --recover 

# # init chain
# testchaind init test --chain-id xion-mainnet-1 --home /Users/donglieu/script/xion/node4
testchaind init test --chain-id xion-mainnet-1

# Change parameter token denominations to stake
cat $HOME/.testchain/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="stake"' > $HOME/.testchain/config/tmp_genesis.json && mv $HOME/.testchain/config/tmp_genesis.json $HOME/.testchain/config/genesis.json
cat $HOME/.testchain/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="stake"' > $HOME/.testchain/config/tmp_genesis.json && mv $HOME/.testchain/config/tmp_genesis.json $HOME/.testchain/config/genesis.json
cat $HOME/.testchain/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="stake"' > $HOME/.testchain/config/tmp_genesis.json && mv $HOME/.testchain/config/tmp_genesis.json $HOME/.testchain/config/genesis.json
cat $HOME/.testchain/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="stake"' > $HOME/.testchain/config/tmp_genesis.json && mv $HOME/.testchain/config/tmp_genesis.json $HOME/.testchain/config/genesis.json

sed -i -E 's|minimum-gas-prices = ""|minimum-gas-prices = "0.0001stake"|g' $HOME/.testchain/config/app.toml

# Allocate genesis accounts (cosmos formatted addresses)
testchaind genesis add-genesis-account cosmos1f7twgcq4ypzg7y24wuywy06xmdet8pc4473tnq 1000000000000stake,100000000uxion --keyring-backend test

# # # # Sign genesis transaction
testchaind genesis gentx val  1000000stake --keyring-backend test --chain-id xion-mainnet-1

# # Collect genesis tx
testchaind genesis collect-gentxs

# # Run this to ensure everything worked and that the genesis file is setup correctly
testchaind genesis validate-genesis

# # Start the node (remove the --pruning=nothing flag if historical queries are not needed)
# # screen -S xionx -t xionx -d -m
# testchaind start

# sleep 7

# testchaind tx bank send $val2 $test2 100000stake  --chain-id realio_3-2 --keyring-backend test --fees 10stake -y #--node tcp://127.0.0.1:26657