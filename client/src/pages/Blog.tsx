import React, { useState, useEffect } from "react";
import PostPreview from "../components/PostPreview";
import { getAllPosts } from "../interfaces/contract_interface";
import { Box } from "@chakra-ui/react";

function Blog() {
  const [posts, setPosts] = useState([{title: "", cid: "", author: ""}])

  function getPosts() {
      getAllPosts().then((r) => (
        setPosts(r)
      ))
  }

  useEffect(() => {
    getPosts()
  }, []);

  return (
    <Box>
      {
        posts.length > 0 && posts.map((p, i) => (
          <PostPreview key={i} title={p.title} />
        ))
      }
    </Box>
  )
  
}

export default Blog;
