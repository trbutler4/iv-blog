
import {
    Box,
    Text,
    LinkBox,
    LinkOverlay
} from '@chakra-ui/react';

interface IPostPreview {
    title: string
}
export default function PostPreview({ title }: IPostPreview) {
    return (
        <LinkBox m={10} background="gray">
            <LinkOverlay href='#/post'>
                {title}    
            </LinkOverlay>
        </LinkBox>
    );
}

