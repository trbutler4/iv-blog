import { Contract, BrowserProvider } from "ethers";
import { get } from "http";

const provider = new BrowserProvider(window.ethereum);

const contractAbi = [{"inputs":[{"internalType":"string","name":"_title","type":"string"},{"internalType":"string","name":"_cid","type":"string"}],"name":"createPost","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getCurrentPostId","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"_postId","type":"uint256"}],"name":"getPost","outputs":[{"components":[{"internalType":"string","name":"title","type":"string"},{"internalType":"string","name":"cid","type":"string"},{"internalType":"address","name":"author","type":"address"}],"internalType":"struct Blog.Post","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"posts","outputs":[{"internalType":"string","name":"title","type":"string"},{"internalType":"string","name":"cid","type":"string"},{"internalType":"address","name":"author","type":"address"}],"stateMutability":"view","type":"function"}]

const contractAddress = "0x3A77966834832939670ABeeEaDBe423696faA113"

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
  const post = await blog.getPost(postId);
  return {
    title: post[0],
    cid: post[1],
    author: post[2]
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
