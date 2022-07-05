import React, {useEffect} from "react";
import DropDown from "../components/Dropdown";
import Checkbox from "../components/Checkbox";
import Button from "../components/Button";
import Link from "next/link";
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
    const [buttonLoading, setButtonLoading] = React.useState(false);

    const ravintolaBoxes = GenerateCheckboxes(ravintolat)
    const [ravintola_lista, lisaaRavintola] = React.useState({
        ravintolat: [],
    })

    const kaupungitBoxes = GenerateCheckboxes(kaupungit)
    const [kaupunki_lista, lisaaKaupunki] = React.useState({
        kaupungit: [],
    })

    //TODO: Saa nama latautumaan silloin, kun asetukset page avataan.

    useEffect(() => {
        if (kaupunki_lista.kaupungit.length > 1) {
            window.localStorage.setItem("varaapoyta_kaupungit", JSON.stringify(kaupunki_lista.kaupungit))
        }
    }, [kaupunki_lista.kaupungit]);

    useEffect(() => {
        if (ravintola_lista.ravintolat.length > 1) {
            window.localStorage.setItem("varaapoyta_ravintolat", JSON.stringify(ravintola_lista.ravintolat))
        }
    }, [ravintola_lista.ravintolat]);

    //TODO: Make this work.
    useEffect(() => {
        // @ts-ignore
        const ravintolat_storage = JSON.parse(window.localStorage.getItem("varaapoyta_ravintolat")) || [];
        // @ts-ignore
        const kaupungit_storage = JSON.parse(window.localStorage.getItem("varaapoyta_kaupungit")) || [];

        lisaaKaupunki({kaupungit: kaupungit_storage})
        lisaaRavintola({ravintolat: ravintolat_storage})
    }, []);

    function setButton() {
        setButtonLoading(true)
        setButtonLoading(false)
    }

    console.log(ravintola_lista.ravintolat)
    console.log(kaupunki_lista.kaupungit)

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