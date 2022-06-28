import Link from 'next/link'
import GitHubIcon from '@mui/icons-material/GitHub';
import LinkedInIcon from '@mui/icons-material/LinkedIn';

function SiteFooter() {
    return (
        <div className={"fixed left-0 bottom-3 w-full text-center"}>
            <p>Created by Rase</p>
            <Link href="https://github.com/52617365">
                <GitHubIcon/>
            </Link>
            <Link href="https://www.linkedin.com/mwlite/in/rasmus-m-4a7195220">
                <LinkedInIcon/>
            </Link>
        </div>
    )
}

export default SiteFooter;