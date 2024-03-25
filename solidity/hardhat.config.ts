import "dotenv/config";
import "@nomiclabs/hardhat-waffle";
import "hardhat-gas-reporter";
import "@typechain/hardhat";
import "hardhat-contract-sizer";
import { task, extendEnvironment } from "hardhat/config";
import {
  HardhatNetworkAccountsUserConfig,
  HardhatUserConfig,
} from "hardhat/types";
import { ethers } from "ethers";
const count = Number(process.env.ACCOUNT_TOTAL || 125); // equal powers
const accountsBalance =
  process.env.ACCOUNT_BALANCE || "10000000000000000000000";
let accounts: HardhatNetworkAccountsUserConfig | undefined = undefined;
if (process.env.MNEMONIC) {
  accounts = {
    mnemonic: process.env.MNEMONIC,
    path: "m/44'/60'/0'/0",
    initialIndex: 0,
    count,
    passphrase: "",
    accountsBalance,
  };
} else if (process.env.PRIVATE_KEY) {
  accounts = process.env.PRIVATE_KEY.split(/\s*,\s*/).map((pv) => ({
    privateKey: pv,
    balance: accountsBalance,
  }));
}

task("accounts", "Prints the list of accounts", async (args, hre) => {
  const accounts = hre.getSigners();

  for (const account of accounts) {
    console.log(
      await account.getAddress(),
      (await account.getBalance()).toString()
    );
  }
});

// You have to export an object to set up your config
// This object can have the following optional entries:
// defaultNetwork, networks, solc, and paths.
// Go to https://buidler.dev/config/ to learn more
const config: HardhatUserConfig = {
  defaultNetwork: "hardhat",
  // This is a sample solc configuration that specifies which version of solc to use
  solidity: {
    compilers: [
      {
        version: "0.8.10",
        settings: {
          optimizer: {
            enabled: true,
          },
        },
      },
      {
        version: "0.8.12",
        settings: {
          optimizer: {
            enabled: true,
          },
        },
      },
    ],
  },
  networks: {
    hardhat: {
      chainId: 420,
      accounts,
      forking: {
        url: "https://rpc.ankr.com/eth_goerli",
        blockNumber: 8218229,
        enabled: false, // turn off because of goerli is dead
      },
      mining: {
        auto: false,
        interval: 2000,
      },
    },
  },
  typechain: {
    outDir: "typechain",
    target: "ethers-v5",
  },
  gasReporter: {
    enabled: true,
  },
  mocha: {
    timeout: process.env.TIMEOUT || 2000000,
  },
};

declare module "hardhat/types/runtime" {
  export interface HardhatRuntimeEnvironment {
    provider: ethers.providers.Web3Provider;
    getSigner: (
      addressOrIndex?: string | number
    ) => ethers.providers.JsonRpcSigner;
    getSigners: (num?: number) => ethers.providers.JsonRpcSigner[];
  }
}

extendEnvironment((hre) => {
  // @ts-ignore
  hre.provider = new ethers.providers.Web3Provider(hre.network.provider);
  hre.getSigners = (num = count) =>
    [...new Array(num)].map((_, i) => hre.provider.getSigner(i));
  hre.getSigner = (addressOrIndex) => hre.provider.getSigner(addressOrIndex);
});

export default config;
