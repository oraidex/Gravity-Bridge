import { ethers } from "hardhat";
import { Overrides } from "@ethersproject/contracts";
import { makeCheckpoint, getSignerAddresses, ZeroAddress } from "./pure";
import { Signer } from "ethers";
import {
  ERC20__factory,
  Gravity__factory,
  TestERC20Custom__factory,
} from "../typechain";

export async function deployContracts(
  gravityId: string = "foo",
  validators: Signer[],
  powers: number[],
  opts?: Overrides
) {
  // enable automining for these tests
  await ethers.provider.send("evm_setAutomine", [true]);
  const [owner] = await ethers.getSigners();
  const testERC20 = await new TestERC20Custom__factory(owner).deploy(
    [await owner.getAddress()],
    opts
  );

  const valAddresses = await getSignerAddresses(validators);

  const gravity = await new Gravity__factory(owner).deploy(
    gravityId,
    valAddresses,
    powers,
    await owner.getAddress(),
    opts
  );

  const checkpoint = makeCheckpoint(
    valAddresses,
    powers,
    0,
    0,
    ZeroAddress,
    gravityId
  );

  return { gravity, testERC20, checkpoint };
}
