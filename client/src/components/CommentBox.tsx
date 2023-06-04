import {
    Box,
    Textarea,
    Button,
    HStack,
    Text
} from '@chakra-ui/react';

function CommentBox() {

    return (
        <Box m={20}> 
            <Textarea placeholder='Add your thoughts' />   
            
            <HStack justify='space-between' align='center'>
                <Text>Address or ENS name here</Text>
                <Button>Comment</Button>
            </HStack>
        </Box>
    )
}

export default CommentBox;
