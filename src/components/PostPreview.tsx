
import {
    Box,
    Text,
    LinkBox,
    LinkOverlay
} from '@chakra-ui/react';

function PostPreview() {
    return (
        <LinkBox m={10} background="gray">
            <LinkOverlay href='#/post'>
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
            </LinkOverlay>
        </LinkBox>
    );
}

export default PostPreview;
