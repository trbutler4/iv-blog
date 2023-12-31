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
            .then((response) => {
                console.log('Content-Type:', response.headers.get('Content-Type'));
                return response.text();
            })
            .then((text) => {
                console.log('Response text:', text);
                setMarkdown(text);
            });
    }, [articleId]);

    return (
        <div className="p-20 prose">
            <Markdown>{markdown}</Markdown>
        </div>
    )
}