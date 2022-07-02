import Link from "next/link";
import Button from "../components/Button"
import SiteFooter from "../components/SiteFooter";
import SiteHead from "../components/SiteHead";

function Varaa() {
    return (
        <>
            <SiteHead/>
            <div className="relative hero min-h-screen bg-base-200">
                <Link href={"/asetukset"}>
                    <a className={"absolute top-0 right-0"}>
                        <Button text={"Muuta asetuksiasi"}/>
                    </a>
                </Link>
                <div className="hero-content text-center">
                    <div className="absolute top-1/3 max-w-md pr-5">
                        {/*TODO: Add different paths to different oauth logins here*/}
                        <h1 className={"text-xl pb-10"}>Raflaamo varaaja</h1>
                        <div className={"grid gap-10"}>
                            <Button text={"Varaa pöytäsi"}/>
                            {/*TODO: Make this button go in top right corner.*/}
                        </div>
                    </div>
                </div>
            </div>
            <SiteFooter/>
        </>
    )
}


export default Varaa;
