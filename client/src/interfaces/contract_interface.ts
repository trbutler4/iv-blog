import { Contract, BrowserProvider } from "ethers";

const provider = new BrowserProvider(window.ethereum);

/*
const contractAbi = ContractData.abi;

const contractAddress = ContractData.address;

const getContractRef = async () => {
  const signer = await provider.getSigner();
  return new Contract(contractAddress, contractAbi, signer);
};
*/
