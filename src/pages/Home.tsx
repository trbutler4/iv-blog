import React from 'react';
import ArticlePreview from '../components/ArticlePreview';

const Home = () => {
  return (
    <div className="flex flex-col p-8 justify-center items-center">
      <ArticlePreview date="2023-12-30" title="Hello World" articlePath="hello-world" />
    </div>
  )
}

export default Home;
