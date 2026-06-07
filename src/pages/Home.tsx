import React from "react";
import ArticlePreview from "../components/ArticlePreview";

const Home = () => {
  return (
    <div className="flex flex-col p-8 justify-center items-center">
      <ArticlePreview
        date="2025-06-06"
        title="Peru"
        articlePath="peru"
      />
    </div>
  );
};

export default Home;
