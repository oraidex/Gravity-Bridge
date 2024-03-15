import { TestERC721A, GravityERC721, TestFakeGravity } from "../typechain";
import { ethers } from "hardhat";
import { getSignerAddresses, ZeroAddress } from "./pure";
import { Signer } from "ethers";

import { deployContracts } from "./index";

export async function deployContractsERC721(
  gravityId: string = "foo",
  validators: Signer[],
  powers: number[]
) {
  const { gravity, testERC20, checkpoint } = await deployContracts(
    gravityId,
    validators,
    powers
  );

  const TestERC721 = await ethers.getContractFactory("TestERC721A");
  const testERC721 = (await TestERC721.deploy()) as TestERC721A;

  const GravityERC721 = await ethers.getContractFactory("GravityERC721");
  const gravityERC721 = (await GravityERC721.deploy(
    gravity.address
  )) as GravityERC721;

  const FakeGravity = await ethers.getContractFactory("TestFakeGravity");
  const fakeGravity = (await FakeGravity.deploy(
    gravityId,
    await getSignerAddresses(validators),
    powers,
    await validators[0].getAddress()
  )) as TestFakeGravity;

  return {
    gravity,
    gravityERC721,
    fakeGravity,
    testERC721,
    testERC20,
    checkpoint,
  };
}
