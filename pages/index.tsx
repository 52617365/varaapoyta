import type {NextPage} from 'next'
import Link from "next/link"
import SiteFooter from "../components/SiteFooter";
import Button from "../components/Button";
import React from "react";

function fetchInfo() {
    const data = {
        operationName: "getRestaurantsByLocation",
        variables: {
            first: 470,
            input: {
                restaurantType: "ALL",
                locationName: "Helsinki",
                feature: {
                    rentableVenues: false
                }
            },
            after: "eyJmIjowLCJnIjp7ImEiOjYwLjE3MTE2LCJvIjoyNC45MzI1OH19"
        },
        query: `fragment Locales on LocalizedString {fi_FI\n }\n\nfragment Restaurant on Restaurant {\n  id\n  name {\n    ...Locales\n    }\n  urlPath {\n    ...Locales\n     }\n    address {\n    municipality {\n      ...Locales\n       }\n        street {\n      ...Locales\n       }\n       zipCode\n     }\n    features {\n    accessible\n     }\n  openingTime {\n    restaurantTime {\n      ranges {\n        start\n        end\n        endNextDay\n         }\n             }\n    kitchenTime {\n      ranges {\n        start\n        end\n        endNextDay\n              }\n             }\n    }\n  links {\n    tableReservationLocalized {\n      ...Locales\n        }\n    homepageLocalized {\n      ...Locales\n          }\n   }\n     \n}\n\nquery getRestaurantsByLocation($first: Int, $after: String, $input: ListRestaurantsByLocationInput!) {\n  listRestaurantsByLocation(first: $first, after: $after, input: $input) {\n    totalCount\n      edges {\n      ...Restaurant\n        }\n     }\n}`
    }

    const headers = {
        method: 'POST',
        mode: 'no-cors',
        headers: {
            "Content-Type": 'application/json',
            "client_id": 'jNAWMvWD9rp637RaR',
        },
        body: JSON.stringify(data)
    }

    console.log(JSON.stringify(data))

    fetch("https://api.raflaamo.fi/query", headers as RequestInit).then(res => console.log(res))
}


// TODO: turn page navigations into regular buttons cuz they load insta instead of having a loading time, only use the loading indicator on buttons that send requests.
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
