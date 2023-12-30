import React from 'react';
import ArticlePreview from '../components/ArticlePreview';

const Home = () => {
    return(
        <div className="flex flex-col p-10"> 
            <ArticlePreview date="2021-10-01" title="First Post" />
        </div>
    )
}

export default Home;
