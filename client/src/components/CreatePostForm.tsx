import {
    Box,
    Text,
    Button,
    Input 
} from '@chakra-ui/react'
import React, { useState } from 'react';
import { createPost } from '../interfaces/contract_interface';

function CreatePostForm() {
    const [title, setTitle] = useState("")
    const [cid, setCid] = useState("")

    function uploadPost(title: string, cid: string) {
        // TODO: call smart contract function with post title and cid
        console.log(`title: ${title}`)
        console.log(`cid: ${cid}`)
        createPost(title, cid).then((r) => (
            console.log(r)
        ))
    }

    return (
        <Box
            padding={20}
            justifyContent="center"
            m={40}
        >
            <Input 
                placeholder='Post Title' 
                size='lg'
                onChange={(e) => setTitle(e.target.value)}
            /> 
            <Input 
                placeholder='CID' 
                size='lg' 
                onChange={(e) => setCid(e.target.value)}
            /> 
            <Button
                onClick={() => uploadPost(title, cid)} 
            >Upload</Button>
        </Box>
    )
}

export default CreatePostForm;
