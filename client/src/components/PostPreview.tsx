
import {
    Text,
    LinkBox,
    LinkOverlay
} from '@chakra-ui/react';

interface IPostPreview {
    title: string
    author: string
}
export default function PostPreview({ title, author }: IPostPreview) {
    return (
        <LinkBox m={5} background="gray.300" marginX={20}>
            <LinkOverlay href='#/post' textAlign="center">
                <Text>{title} -- {author}</Text>
            </LinkOverlay>
        </LinkBox>
    );
}

