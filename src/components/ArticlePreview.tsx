import { Link } from "react-router-dom";

export default function ArticlePreview({
    date, 
    title, 
    articlePath,
}: {
    date: string;
    title: string;
    articlePath: string;
}) {

    function articleRedirect() {
        console.log(`Redirecting to ${articlePath}`)
    }

    return (
        <div className="py-2" >
            <p className="text-xs">{date}</p>
            <Link 
                className="text-xl hover:underline" 
                to={`articles/${articlePath}`}
            >
                {title}
            </Link>
        </div>
    );
}

