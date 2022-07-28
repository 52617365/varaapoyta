import type {NextPage} from 'next'
import Link from "next/link"
import SiteFooter from "../components/SiteFooter";
import Button from "../components/Button";
import React from "react";

const Home: NextPage = () => {
    const [buttonLoading, setButtonLoading] = React.useState(false);

    function fetchInformation(event: any) {
        setButtonLoading(true)
        // TODO: fetchInfo should set buttonLoading to false when fetched.
        fetchInfo()
    }

    return (
        <>
            <div className="hero min-h-screen bg-base-200">
                <div className="hero-content text-center">
                    <div className="max-w-md pr-5">
                        <h1 className="text-5xl font-bold">Moikka,</h1>
                        <p className="py-6">Aloita pöydänvaraus kirjoittamalla kaupunkisi kenttään.</p>
                        <Link href={"/varaa"}>
                            <a className={"absolute top-0 right-0"}>
                                <button className="btn">Varaa</button>
                            </a>
                        </Link>
                        <Link href={"/asetukset"}>
                            <a className={"absolute top-0 right-20"}>
                                <button className="btn">Asetukset</button>
                            </a>
                        </Link>
                        <div className={"pb-3"}>
                            <input type="text" placeholder="Kaupunki" className="input w-full max-w-xs"/>
                        </div>
                        <div>
                            <Button text="Hae ravintolat" setButton={fetchInformation} buttonLoading={buttonLoading}/>
                        </div>
                    </div>
                </div>
            </div>
            <SiteFooter/>
        </>
    )
}
export default Home
