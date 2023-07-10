import { Contract, BrowserProvider } from "ethers";
import { get } from "http";
import contractInfo from "../contractInfo.json";

const provider = new BrowserProvider(window.ethereum);

// TODO: update these to manage better based on dev/prod 
const contractAbi = contractInfo.abi;
const contractAddress = "0x65eaDd42362d714cA0Dac01D0856bCDD40749c95"

async function getContractRef() {
  const signer = await provider.getSigner();
  return new Contract(contractAddress, contractAbi, signer);
};

export async function createPost(title: string, cid: string, url: string) {
  const blog = await getContractRef();
  console.log("creating post...")
  return await blog.createPost(title, cid, url);
}

export async function getPostById(postId: number) {
  const blog = await getContractRef();
  console.log(`getting post ${postId}...`);
  const post = await blog.getPost(postId);
  return {
    title: post[0],
    cid: post[1],
    url: post[2],
    author: post[3]
  }
}

export async function getCurrentPostId() {
  const blog = await getContractRef();
  const postId = await blog.getCurrentPostId();
  return postId;
  
}

export async function getAllPosts() {
  const blog = await getContractRef();
  const currentPostId = await blog.getCurrentPostId();
  
  let posts = [];
  for (let i = 0; i < currentPostId; i++) {
    const p = await blog.getPost(i);
    posts.push(p);
  }

  return posts;
}
