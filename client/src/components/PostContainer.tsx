import {
    Box
} from '@chakra-ui/react';
import { ReactMarkdown } from "react-markdown/lib/react-markdown"
import { useEffect, useState } from 'react';

function PostContainer() {
    const [postContent, setPostContent] = useState("# POST CONTENT");

    async function retrievePostFromIPFS(cid: string) {
        try {
            const response = await fetch('https://bafybeien4wqtatzszfrdpwcy5c56tbefxauzb7og7v2oa7gpfvic7nu24q.ipfs.w3s.link/example0.md');
            const content = await response.text();
            setPostContent(content);
        } catch (e) {
            console.log('Error retrieving post: ', e);
        }
    }

    useEffect(() => {
        const cid = "bafybeien4wqtatzszfrdpwcy5c56tbefxauzb7og7v2oa7gpfvic7nu24q" // TODO: pass this via the link to the post
        retrievePostFromIPFS(cid);
    })

    return (
        <Box m={20}> 
            <ReactMarkdown> 
                {postContent}
           </ReactMarkdown>
        </Box>
    )
}

export default PostContainer;
