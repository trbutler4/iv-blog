import React, { useState, useEffect } from "react";
import PostPreview from "../components/PostPreview";
import { getAllPosts } from "../interfaces/contract_interface";
import { Box, Text } from "@chakra-ui/react";

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
    <Box padding={20}>
        <Text 
            fontSize="3xl"
            textAlign="left"
            marginLeft={20}
        >
            Posts</Text>
        <Box>
          {
            posts.length > 0 && posts.map((p, i) => (
              <PostPreview key={i} title={p.title} author={p.author}/>
            ))
          }
        </Box>
    </Box>
  )
  
}

export default Blog;
