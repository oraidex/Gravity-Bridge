import chai from "chai";
import { ethers } from "hardhat";
import { solidity } from "ethereum-waffle";

import { deployContracts } from "../test-utils";
import {
  getSignerAddresses,
  makeCheckpoint,
  signHash,
  makeTxBatchHash,
  examplePowers,
} from "../test-utils/pure";

chai.use(solidity);
const { expect } = chai;

async function runTest(opts: {}) {
  // Prep and deploy contract
  // ========================
  const signers = await ethers.getSigners();
  const gravityId = ethers.utils.formatBytes32String("foo");
  // This is the power distribution on the Cosmos hub as of 7/14/2020
  let powers = examplePowers();
  let validators = signers.slice(0, powers.length);
  const {
    gravity,
    testERC20,
    checkpoint: deployCheckpoint,
  } = await deployContracts(gravityId, validators, powers);

  // Transfer out to Cosmos, locking coins
  // =====================================
  await testERC20.approve(gravity.address, 1000);
  await expect(
    gravity.sendToCosmos(
      testERC20.address,
      ethers.utils.formatBytes32String("myCosmosAddress"),
      1000
    )
  )
    .to.emit(gravity, "SendToCosmosEvent")
    .withArgs(
      testERC20.address,
      await signers[0].getAddress(),
      ethers.utils.formatBytes32String("myCosmosAddress"),
      1000,
      2
    );

  expect(await testERC20.balanceOf(gravity.address)).to.equal(1000);
  expect(await gravity.state_lastEventNonce()).to.equal(2);

  // Do it again
  // =====================================
  await testERC20.approve(gravity.address, 1000);
  await expect(
    gravity.sendToCosmos(
      testERC20.address,
      ethers.utils.formatBytes32String("myCosmosAddress"),
      1000
    )
  )
    .to.emit(gravity, "SendToCosmosEvent")
    .withArgs(
      testERC20.address,
      await signers[0].getAddress(),
      ethers.utils.formatBytes32String("myCosmosAddress"),
      1000,
      3
    );

  expect(await testERC20.balanceOf(gravity.address)).to.equal(2000);
  expect(await gravity.state_lastEventNonce()).to.equal(3);
}

describe("sendToCosmos tests", function () {
  it("works right", async function () {
    await runTest({});
  });
});
