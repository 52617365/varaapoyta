import Link from "next/link";
import Button from "../components/Button"
import SiteFooter from "../components/SiteFooter";
import {useSWRConfig} from "swr";

function Varaa() {
    const {cache, mutate, ...extraConfig} = useSWRConfig()
    console.log(cache.get("ravintolat"))
    console.log(cache.get("kaupungit"))
    return (
        <>
            <div className="relative hero min-h-screen bg-base-200">
                <Link href={"/asetukset"}>
                    <a className={"absolute top-0 right-0"}>
                        <Button text={"Muuta asetuksiasi"}/>
                    </a>
                </Link>
                <div className="hero-content text-center">
                    <div className="absolute top-1/3 max-w-md pr-5">
                        {/*TODO: Add different paths to different oauth logins here*/}
                        <h1 className={"pb-10 text-xl"}>Raflaamo varaaja</h1>
                        <div className={"grid gap-10 w-full"}>
                            <Button text={"Varaa pöytäsi"}/>
                        </div>
                    </div>
                </div>
            </div>
            <SiteFooter/>
        </>
    )
}


export default Varaa;
