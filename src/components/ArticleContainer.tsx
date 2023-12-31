import { useEffect, useState } from "react";

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
        <div>
            <p>Article Container</p>
            <p>{articleId}</p>
            <p>{markdown}</p>
        </div>
    )
}