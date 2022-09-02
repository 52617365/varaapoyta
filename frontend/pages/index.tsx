import type {NextPage} from 'next'
import Link from "next/link"
import Button from "../components/Button";
import React, {useState} from "react";
import Card from "../components/Card";
import api_response from "../interfaces/api_response_interface";

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
                        <p className="py-6">Aloita pöytävaraus kirjoittamalla kaupunkisi kenttään.</p>
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
                        <div className='pt-5'>
                            <div className="grid grid-rows-4 grid-flow-col gap-4">
                                {ravintolatApista.map((ravintola: api_response) => {
                                    return (
                                        // Storing the id from the reservation page url so its easy to reuse when in V2 we have reservation too. 
                                        <div key={ravintola.links.tableReservationLocalizedId} className='pb-5'>
                                            <Card texts={ravintola}/>
                                        </div>
                                    )
                                })}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </>
    )
}
export default Home
