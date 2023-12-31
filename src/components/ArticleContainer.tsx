import { useEffect, useState } from "react";
import Markdown from "react-markdown";

export default function ArticleContainer({
    articleId,
}: {
    articleId: string;
}) {    
    const [markdown, setMarkdown] = useState<string>("")

    useEffect(() => {
        fetch(`/articles/${articleId}.md`)
            .then((response) => response.text())
            .then((text) => setMarkdown(text))
    }, [articleId]);


    return (
        <div className="p-20">
            <Markdown>{markdown}</Markdown>
        </div>
    )
}