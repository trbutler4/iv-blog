import { Contract, BrowserProvider } from "ethers";
import { get } from "http";

const provider = new BrowserProvider(window.ethereum);

const contractAbi = [{"inputs":[{"internalType":"string","name":"_title","type":"string"},{"internalType":"string","name":"_cid","type":"string"}],"name":"createPost","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getCurrentPostId","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"_postId","type":"uint256"}],"name":"getPost","outputs":[{"components":[{"internalType":"string","name":"title","type":"string"},{"internalType":"string","name":"cid","type":"string"},{"internalType":"address","name":"author","type":"address"}],"internalType":"struct Blog.Post","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"posts","outputs":[{"internalType":"string","name":"title","type":"string"},{"internalType":"string","name":"cid","type":"string"},{"internalType":"address","name":"author","type":"address"}],"stateMutability":"view","type":"function"}];

const contractAddress = "0xde5671D0B5c63689a8DAaC87bAF9796f675D52E7"

async function getContractRef() {
  const signer = await provider.getSigner();
  return new Contract(contractAddress, contractAbi, signer);
};

export async function createPost(title: string, cid: string) {
  const blog = await getContractRef();
  console.log("creating post...")
  return await blog.createPost(title, cid);
}

export async function getPostById(postId: number) {
  const blog = await getContractRef();
  console.log(`getting post ${postId}...`);
  return await blog.getPost(postId);
}