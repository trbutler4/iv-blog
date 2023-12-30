import React from 'react';
import ArticlePreview from '../components/ArticlePreview';

const Home = () => {
    return(
        <div className="flex flex-col p-20"> 
            <ArticlePreview date="2023-12-30" title="Hello World" />
        </div>
    )
}

export default Home;
