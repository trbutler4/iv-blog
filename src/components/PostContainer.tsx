import {
    Box
} from '@chakra-ui/react';
import { ReactMarkdown } from "react-markdown/lib/react-markdown"

function PostContainer() {
    return (
        <Box m={20}> 
            <ReactMarkdown>  
                # THIS IS A POST 
                here is a bunch of text for the post. this is in markdown
                right there in react. very cool! super dope!  
           </ReactMarkdown>
        </Box>
    )
}

export default PostContainer;
