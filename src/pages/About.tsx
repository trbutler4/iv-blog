import React from 'react';
import {
    Text,
    Box
} from '@chakra-ui/react';

function About() {
    return(
        <Box m={10}>
            <Text>
                This is entended to be a place to share my thoughts and ideas in the hope
                that someone will find them useful. I will be wrong often, and I hope that
                I will be corrected. I do not claim to be an expert in any of the topics
                that I am going to write about.
            </Text>
        </Box>
    )
}

export default About;
