import type {NextPage} from 'next'
import Link from "next/link"
import SiteFooter from "../components/SiteFooter";
import Button from "../components/Button";
import React, {useState} from "react";
import {webpack} from "next/dist/compiled/webpack/webpack";

interface string_field {
    fi_FI: string
}

interface ranges_contents {
    end: string;
    start: string;
}

interface ranges {
    ranges: Array<ranges_contents>
}

interface opening_time {
    kitchenTime: ranges;
    restaurantTime: ranges;
}
interface address {
    municipality: string_field;
    street: string_field;
    zipcode: string
}

interface api_response {
    id: number;
    address: address
    available_time_slots: Array<string>;
    name: string_field;
    openingTime: opening_time;
    urlPath: string_field;
}

const Home: NextPage = () => {
    const [buttonLoading, setButtonLoading] = React.useState(false);
    const [ravintolatApista, setRavintolat] = React.useState<api_response[]>([]);

    const fetchInformation = async (city: string) => {
        if (buttonLoading || city == "") {
            return
        }
        setButtonLoading(true)
        try {
            const url = `http://localhost:10000/tables/${city}/1`
            const response = await fetch(url)
            const parsed_response = await response.json()
            console.log(parsed_response)
            setRavintolat(parsed_response)
            setButtonLoading(false)
        } catch (e) {
            setButtonLoading(false)
            console.log("Error fetching endpoint.")
        }
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
                            <input type="text" placeholder="Kaupunki" className="input w-full max-w-xs"
                                   onChange={handleKaupunki}/>
                        </div>
                        <div>
                            <Button text="Hae ravintolat" setButton={fetchInformation} buttonLoading={buttonLoading}
                                    textfield_text={kaupunki}/>
                        </div>
                        <div>
                            {ravintolatApista.map((ravintola: api_response) => {
                                return (
                                    <>
                                        <h1>{ravintola.id}</h1>
                                        <h1>{ravintola.address.municipality.fi_FI}</h1>
                                        <h1>{ravintola.address.street.fi_FI}</h1>
                                        <h1>{ravintola.address.zipcode}</h1>
                                    </>
                                )
                            })}
                        </div>
                    </div>
                </div>
            </div>
            <SiteFooter/>
        </>
    )
}
export default Home
