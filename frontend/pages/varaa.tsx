import Link from "next/link";
import SiteFooter from "../components/SiteFooter";
import React from "react";

function Varaa() {
    return (
        <>
            <div className="relative hero min-h-screen bg-base-200">
                <Link href={"/asetukset"}>
                    <a className={"absolute top-0 right-0"}>
			<button className="btn">Muuta asetuksiasi</button>
                    </a>
                </Link>
                <div className="hero-content text-center">
                    <div className="absolute top-1/3 max-w-md pr-5">
                        <h1 className={"pb-10 text-xl"}>Raflaamo varaaja</h1>
                        <div className={"grid gap-10 w-full"}>
                                <button className="btn">Varaa pöytäsi</button>
                        </div>
                    </div>
                </div>
            </div>
            <SiteFooter/>
        </>
    )
}


export default Varaa;
