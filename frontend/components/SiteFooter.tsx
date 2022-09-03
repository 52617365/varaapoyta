import Link from "next/link";
import GitHubIcon from "@mui/icons-material/GitHub";
import LinkedInIcon from "@mui/icons-material/LinkedIn";

function SiteFooter() {
  return (
    <footer className="footer footer-center p-2 text-base-content bottom-0 fixed bg-base-200">
      <div>
        <p>Copyright © 2022 - All right reserved</p>
        <p>Created by Rase</p>
        <div>
          <Link href="https://github.com/52617365">
            <GitHubIcon />
          </Link>
          {/* TODO: fix the link cuz its old and expired  */}
          <Link href="https://www.linkedin.com/mwlite/in/rasmus-m-4a7195220">
            <LinkedInIcon />
          </Link>
        </div>
      </div>
    </footer>
  );
}

export default SiteFooter;
