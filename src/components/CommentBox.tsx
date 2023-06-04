import {
    Box,
    Input,
    Button
} from '@chakra-ui/react';

function CommentBox() {

    return (
        <Box m={20}> 
            <Input placeholder='Add your thoughts' />   
           
            <Box>
                <Button>Comment</Button>
            </Box>
        </Box>
    )
}

export default CommentBox;
