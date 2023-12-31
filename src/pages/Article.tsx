import ArticleContainer from "../components/ArticleContainer"
import { useParams } from "react-router-dom"

export default function Article() {
    // get route params 
    const { articleId } = useParams();
    console.log(articleId);

    if (!articleId) return (
        <div>
            <p>Article not found</p>
        </div>
    )

    return (
        <div className="flex flex-col justify-center items-center">
            <ArticleContainer articleId={articleId}/>
        </div>
    )
}