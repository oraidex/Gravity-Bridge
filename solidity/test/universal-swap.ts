import { expect } from "chai";
import { ethers } from "hardhat";
import {
  CWSimulateApp,
  GenericError,
  IbcOrder,
  IbcPacket,
  SimulateCosmWasmClient,
} from "@oraichain/cw-simulate";
import { coins, coin } from "@cosmjs/stargate";
import {
  CwIcs20LatestClient,
  Cw20BaseClient,
} from "@oraichain/common-contracts-sdk";
import * as commonArtifacts from "@oraichain/common-contracts-build";
import bech32 from "bech32";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { toBinary } from "@cosmjs/cosmwasm-stargate";

import { deployContracts } from "../test-utils";
import { examplePowers } from "../test-utils/pure";
import { Gravity, TestERC20Custom } from "../typechain";

const senderAddress = "orai1xzmgjjlz7kacgkpxk5gn6lqa0dvavg8r9ng2vu";
const oraibSenderAddress = bech32.encode(
  "oraib",
  bech32.decode(senderAddress).words
);
const gravityId = ethers.utils.formatBytes32String("oraib");

describe("sendToCosmos with IBC Wasm tests", function () {
  let oraibridgeChain: CWSimulateApp;
  // oraichain support cosmwasm
  let oraiClient: SimulateCosmWasmClient;
  let oraiPort: string;
  let oraiIbcDenom: string;
  let ics20: CwIcs20LatestClient;
  let erc20: TestERC20Custom;
  let cw20: Cw20BaseClient;
  let gravity: Gravity;
  let owner: SignerWithAddress;
  let channel = "channel-0";

  beforeEach(async () => {
    oraibridgeChain = new CWSimulateApp({
      chainId: "oraibridge-subnet-2",
      bech32Prefix: "oraib",
    });

    oraiClient = new SimulateCosmWasmClient({
      chainId: "Oraichain",
      bech32Prefix: "orai",
    });

    cw20 = new Cw20BaseClient(
      oraiClient,
      senderAddress,
      (
        await commonArtifacts.deployContract(
          oraiClient,
          senderAddress,
          {
            decimals: 6,
            symbol: "ORAI",
            name: "ERC20 token",
            mint: { minter: senderAddress },
            initial_balances: [
              { address: senderAddress, amount: "1000000000" },
            ],
          },
          "cw20-base"
        )
      ).contractAddress
    );

    ics20 = new CwIcs20LatestClient(
      oraiClient,
      senderAddress,
      (
        await commonArtifacts.deployContract(
          oraiClient,
          senderAddress,
          {
            allowlist: [],
            default_timeout: 3600,
            gov_contract: senderAddress,
            swap_router_contract: "placeholder",
          },
          "cw-ics20-latest"
        )
      ).contractAddress
    );

    oraiPort = "wasm." + ics20.contractAddress;
    // topup
    oraiClient.app.bank.setBalance(
      ics20.contractAddress,
      coins("10000000000000", "orai")
    );

    // Prep and deploy contract
    // ========================
    const signers = await ethers.getSigners();
    // This is the power distribution on the Cosmos hub as of 7/14/2020
    let powers = examplePowers();
    let validators = signers.slice(0, powers.length);
    const res = await deployContracts(gravityId, validators, powers);
    erc20 = res.testERC20;
    gravity = res.gravity;
    owner = signers[0];
    oraiIbcDenom = oraibridgeChain.bech32Prefix + erc20.address;

    // mapping pair from evm erc20 to cosmos cw20
    await ics20.updateMappingPair({
      localAssetInfo: {
        token: {
          contract_addr: cw20.contractAddress,
        },
      },
      localAssetInfoDecimals: 6,
      denom: oraiIbcDenom,
      remoteDecimals: 6,
      localChannelId: channel,
    });

    // init ibc channel between two chains
    oraiClient.app.ibc.relay(
      channel,
      oraiPort,
      channel,
      "transfer",
      oraibridgeChain
    );

    await oraibridgeChain.ibc.sendChannelOpen({
      open_init: {
        channel: {
          counterparty_endpoint: {
            port_id: oraiPort,
            channel_id: channel,
          },
          endpoint: {
            port_id: "transfer",
            channel_id: channel,
          },
          order: IbcOrder.Unordered,
          version: "ics20-1",
          connection_id: "connection-0",
        },
      },
    });

    await oraibridgeChain.ibc.sendChannelConnect({
      open_ack: {
        channel: {
          counterparty_endpoint: {
            port_id: oraiPort,
            channel_id: channel,
          },
          endpoint: {
            port_id: "transfer",
            channel_id: channel,
          },
          order: IbcOrder.Unordered,
          version: "ics20-1",
          connection_id: "connection-0",
        },
        counterparty_version: "ics20-1",
      },
    });

    // handle IBC callback
    oraibridgeChain.ibc.addMiddleWare(async (msg, app) => {
      // check if memo is json
      try {
        const { wasm } = JSON.parse(msg.data.memo);
        if (wasm)
          await oraiClient.app.wasm.executeContract(
            senderAddress,
            [],
            wasm.contract,
            wasm.msg
          );
      } catch {}
    });

    // handle evm event callback
    gravity.on(gravity.filters.SendToCosmosEvent(), async (...args) => {
      const [erc20Addr, , destination, amount] = args;
      const denom = oraibridgeChain.bech32Prefix + erc20Addr;
      expect(await erc20.balanceOf(gravity.address)).to.equal(amount);
      expect(await gravity.state_lastEventNonce()).to.equal(2);

      // Relay to Oraichain
      console.log("destination", destination);
      console.dir(await ics20.channel({ id: channel }), { depth: null });

      await oraibridgeChain.ibc.sendPacketReceive({
        packet: {
          data: toBinary({
            amount: amount.toString(),
            denom,
            receiver: destination,
            sender: oraibSenderAddress,
            memo: "",
          }),
          src: {
            port_id: "transfer",
            channel_id: channel,
          },
          dest: {
            port_id: oraiPort,
            channel_id: channel,
          },
          sequence: 27,
          timeout: {
            block: {
              revision: 1,
              height: 12345678,
            },
          },
        },
        relayer: oraibSenderAddress,
      });

      // check Orai erc20 token sent to Oraichain via IBC Wasm channel
      console.dir(await ics20.channel({ id: channel }), { depth: null });
    });
  });

  it("send token from evm to oraichain via oraibridge", async function () {
    // Transfer out to Cosmos, locking coins
    // =====================================

    const amount = BigInt(100000000);
    await erc20.approve(gravity.address, amount);

    await gravity.sendToCosmos(erc20.address, senderAddress, amount);

    // wait 5s due to hardhat pooling of 4s
    await new Promise((resolve) => setTimeout(resolve, 5000));
  });
});
