import {
    Box
} from '@chakra-ui/react';
import { ReactMarkdown } from "react-markdown/lib/react-markdown"
import { useState } from 'react';
import * as fs from 'fs';

function PostContainer() {
    const [postContent, setPostContent] = useState("# POST CONTENT");

    return (
        <Box m={20}> 
            <ReactMarkdown> 
                {postContent}
           </ReactMarkdown>
        </Box>
    )
}

export default PostContainer;
