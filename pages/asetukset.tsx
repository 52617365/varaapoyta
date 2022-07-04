import React from "react";
import DropDown from "../components/Dropdown";
import Checkbox from "../components/Checkbox";
import Button from "../components/Button";
import Link from "next/link";
import {useSWRConfig} from "swr";
import SaveButton from "../components/SaveButton"

function GenerateCheckboxes(items: string[]) {
    return (
        items.map((item) => {
            return (
                <Checkbox key={item} name={item}/>
            )
        }))
}

function RavintolaHandler(e: any, settings: any, setSettings: any) {
    // Destructuring
    const {name, checked} = e.target;
    const {ravintolat} = settings;

    // Case 1 : The user checks the box
    if (checked) {
        setSettings({
            // @ts-ignore
            ravintolat: [...ravintolat, name],
        });
    }
    // Case 2  : The user unchecks the box
    else {
        setSettings({
            ravintolat: ravintolat.filter((e: string) => e !== name),
        });
    }
}

function KaupunkiHandler(e: any, settings: any, setSettings: any) {
    // Destructuring
    const {name, checked} = e.target;
    const {kaupungit} = settings;

    // Case 1 : The user checks the box
    if (checked) {
        setSettings({
            // @ts-ignore
            kaupungit: [...kaupungit, name],
        });
    }
    // Case 2  : The user unchecks the box
    else {
        setSettings({
            kaupungit: kaupungit.filter((e: string) => e !== name),
        });
    }
}

function Asetukset({ravintolat, kaupungit}: { ravintolat: string[], kaupungit: string[] }) {
    const {cache, mutate, ...extraConfig} = useSWRConfig()
    const [buttonLoading, setButtonLoading] = React.useState(false);

    const ravintolaBoxes = GenerateCheckboxes(ravintolat)
    const [ravintola_lista, lisaaRavintola] = React.useState({
        ravintolat: [],
    })

    const kaupungitBoxes = GenerateCheckboxes(kaupungit)
    const [kaupunki_lista, lisaaKaupunki] = React.useState({
        kaupungit: [],
    })

    function setButton() {
        setButtonLoading(true)
        cache.set("ravintolat", ravintola_lista)
        cache.set("kaupungit", kaupunki_lista)
        setButtonLoading(false)
    }

    console.log(cache.get("ravintolat"))
    console.log(cache.get("kaupungit"))

    return (
        <>
            <div className="relative hero min-h-screen bg-base-200">
                <Link href={"/varaa"}>
                    <a className={"absolute top-0 right-0"}>
                        <Button text={"Takaisin varaamaan"}/>
                    </a>
                </Link>
                <div className="hero-content text-center">
                    <div className="absolute top-1/3 max-w-md pr-5">
                        {/*TODO: Add different paths to different oauth logins here*/}
                        <h1 className={"pb-10 text-xl"}>Asetukset</h1>
                        <div className={"grid gap-5 w-full"}>
                            <DropDown
                                onChange={val => RavintolaHandler(val, ravintola_lista, lisaaRavintola)}
                                items={ravintolaBoxes}
                                name={"Ravintolat"}
                            />
                            <DropDown
                                onChange={val => KaupunkiHandler(val, kaupunki_lista, lisaaKaupunki)}
                                items={kaupungitBoxes}
                                name={"Kaupungit"}
                            />
                            <div onClick={setButton}>
                                <SaveButton text={"Tallenna asetukset"} buttonLoading={buttonLoading}
                                            setButton={setButton}/>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </>
    )
}

// TODO: Add other data in here too.
export async function getStaticProps() {
    // TODO: Hae ravintola nimet jostain ja anna ne main componenttiin tasta.
    const ravintolat = ["restaurant1", "restaurant2", "restaurant3"];
    const kaupungit = ["Rovaniemi", "Helsinki", "Tampere"]
    return {
        props: {
            ravintolat: ravintolat,
            kaupungit: kaupungit
        },
    }
}

export default Asetukset;