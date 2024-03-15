import {
  TestERC721A__factory,
  GravityERC721__factory,
  TestFakeGravity__factory,
} from "../typechain";
import { ethers } from "hardhat";
import { getSignerAddresses } from "./pure";
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

  const [owner] = await ethers.getSigners();

  const testERC721 = await new TestERC721A__factory(owner).deploy();

  const gravityERC721 = await new GravityERC721__factory(owner).deploy(
    gravity.address
  );

  const fakeGravity = await new TestFakeGravity__factory(owner).deploy(
    gravityId,
    await getSignerAddresses(validators),
    powers,
    await validators[0].getAddress()
  );

  return {
    gravity,
    gravityERC721,
    fakeGravity,
    testERC721,
    testERC20,
    checkpoint,
  };
}
