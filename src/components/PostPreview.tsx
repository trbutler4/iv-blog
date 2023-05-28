
import {
    Box,
    Text
} from '@chakra-ui/react';

function PostPreview() {
    return (
        <Box m={10} background="gray">

            <Box display="flex">
                <Text>Post Title</Text>
                <Text marginLeft="auto">Post Date</Text>
            </Box>

            <Box p={10}>
                <Text>
                    This is going to be a preview of the blog post contents.
                    It will be a short preview of the contents of the blog post.
                </Text>
            </Box>
        </Box>
    );
}

export default PostPreview;
