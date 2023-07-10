import {
    Box,
    Text,
    Button,
    Input 
} from '@chakra-ui/react'
import React, { useState } from 'react';
import { createPost } from '../interfaces/contract_interface';

function UploadPostForm() {
    const [title, setTitle] = useState("")
    const [cid, setCid] = useState("")
    const [url, setUrl] = useState("")

    function uploadPost(title: string, cid: string, url: string) {
        // TODO: call smart contract function with post title and cid
        console.log(`title: ${title}`)
        console.log(`cid: ${cid}`)
        console.log(`url: ${url}`)
        createPost(title, cid, url).then((r) => (
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
            <Input 
                placeholder='URL' 
                size='lg' 
                onChange={(e) => setCid(e.target.value)}
            /> 
            <Button
                onClick={() => uploadPost(title, cid, url)} 
            >Upload</Button>
        </Box>
    )
}

export default UploadPostForm;
