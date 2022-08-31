import type {NextPage} from 'next'
import Link from "next/link"
import SiteFooter from "../components/SiteFooter";
import Button from "../components/Button";
import React, {useState} from "react";

const Home: NextPage = () => {
    const [buttonLoading, setButtonLoading] = React.useState(false);
    const [ravintolatApista, setRavintolat] = React.useState(null);

    const fetchInformation = async (city: string) => {
        if (buttonLoading || city == "") {
            return
        }
        setButtonLoading(true)
        try {
            const url = `http://localhost:10000/tables/${city}/1`
            const response = await fetch(url)
            const response_json = await response.json()
            console.log(response_json)
            setRavintolat(response_json)
            setButtonLoading(false)
        } catch (e) {
            console.log("Error fetching endpoint.")
        }
        // TODO: fetchInfo should set buttonLoading to false when fetched.
        setButtonLoading(false)
    }
    // kaupunki is used in get query to endpoint
    const [kaupunki, asetaKaupunki] = useState('');
    const handleKaupunki = (event: React.ChangeEvent<HTMLInputElement>) => {
        asetaKaupunki(event.target.value);
    };
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
                            <input type="text" placeholder="Kaupunki" className="input w-full max-w-xs" onChange={handleKaupunki}/>
                        </div>
                        <div>
                            <Button text="Hae ravintolat" setButton={fetchInformation} buttonLoading={buttonLoading} textfield_text={kaupunki}/>
                        </div>
                    </div>
                </div>
            </div>
            <SiteFooter/>
        </>
    )
}
export default Home
