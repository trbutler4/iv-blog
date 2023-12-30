

export default function ArticlePreview({
    date, 
    title, 
}: {
    date: string;
    title: string;
}) {
    return (
        <div>
            <p className="text-xs">{date}</p>
            <p className="text-xl">{title}</p>
        </div>
    );
}

