import React, {useEffect} from "react";
import DropDown from "../components/Dropdown";
import Checkbox from "../components/Checkbox";
import Button from "../components/Button";
import Link from "next/link";

// TODO: Fix checkboxes not updating on first page load, instead they generate on refresh.
function GenerateCheckboxes(items: string[], storage_items: string[]) {
    // At this time, storage_items will contain all the items from the local storage. So here we check if it's in the local storage, we mark the checkbox as checked, else its unchecked.
    return (
        items.map((item) => {
            {
                if (storage_items.includes(item)) {
                    return (
                        <Checkbox key={item} name={item} checked={true}/>
                    )
                } else {
                    return (
                        <Checkbox key={item} name={item} checked={false}/>
                    )
                }
            }
        }))
}

function RavintolaHandler(e: any, settings: any, setSettings: any) {
    // Destructuring
    const {name, checked} = e.target;
    const {ravintolat} = settings;

    // Case 1: The user checks the box
    if (checked) {
        // checking if the state alrdy has the value we're about to set to avoid duplicates.
        if (!settings.ravintolat.includes(name)) {
            setSettings({
                ravintolat: [...ravintolat, name],
            });
        }
    }
    // Case 2: The user unchecks the box
    else {
        setSettings({
            ravintolat: ravintolat.filter((e: string) => e !== name),
        });
    }
}

// tama on kaupunkeja varten.
function KaupunkiHandler(e: any, settings: any, setSettings: any) {
    // Destructuring
    const {name, checked} = e.target;
    const {kaupungit} = settings;

    // Case 1: The user checks the box
    if (checked) {
        // checking if the state alrdy has the value we're about to set to avoid duplicates.
        if (!settings.kaupungit.includes(name)) {
            setSettings({
                kaupungit: [...kaupungit, name],
            });
        }
    }
    // Case 2: The user unchecks the box
    else {
        setSettings({
            kaupungit: kaupungit.filter((e: string) => e !== name),
        });
    }
}

// TODO: Laita local storagesta haetut setit checkboxeihin ja checkaa ne checkboxit.
function Asetukset({ravintolat, kaupungit}: { ravintolat: string[], kaupungit: string[] }) {
    const [ravintola_lista, lisaaRavintola] = React.useState({
        ravintolat: [],
    })
    const ravintolaBoxes = GenerateCheckboxes(ravintolat, ravintola_lista.ravintolat)

    const [kaupunki_lista, lisaaKaupunki] = React.useState({
        kaupungit: [],
    })
    const kaupungitBoxes = GenerateCheckboxes(kaupungit, kaupunki_lista.kaupungit)


    useEffect(() => {
        if (kaupunki_lista.kaupungit.length > 1) {
            window.localStorage.setItem("varaapoyta_kaupungit", JSON.stringify(kaupunki_lista.kaupungit))
            const items = window.localStorage.getItem("varaapoyta_kaupungit")
            if (items != null) {
                const parsed_items = JSON.parse(items)
                const deduplicated_items = parsed_items.filter((c: any, index: any) => {
                    return parsed_items.indexOf(c) === index;
                });
                window.localStorage.setItem("varaapoyta_kaupungit", JSON.stringify(deduplicated_items))
            }
        }
    }, [kaupunki_lista]);

    useEffect(() => {
        if (ravintola_lista.ravintolat.length > 1) {
            window.localStorage.setItem("varaapoyta_ravintolat", JSON.stringify(ravintola_lista.ravintolat))
            const items = window.localStorage.getItem("varaapoyta_ravintolat")
            if (items != null) {
                const parsed_items = JSON.parse(items)
                const deduplicated_items = parsed_items.filter((c: any, index: any) => {
                    return parsed_items.indexOf(c) === index;
                });
                window.localStorage.setItem("varaapoyta_ravintolat", JSON.stringify(deduplicated_items))
            }
        }
    }, [ravintola_lista]);

    useEffect(() => {
        const ravintolat_storage = JSON.parse(window.localStorage.getItem("varaapoyta_ravintolat") as string) || [];
        lisaaRavintola({ravintolat: ravintolat_storage});

        const kaupunki_storage = JSON.parse(window.localStorage.getItem("varaapoyta_kaupungit") as string) || [];
        lisaaKaupunki({kaupungit: kaupunki_storage});
    }, []);

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
                        <h1 className={"pb-10 text-xl"}>Asetukset</h1>
                        <div className={"grid gap-5 w-full"}>
                            {/*TODO: Add the dropdown here directly because else it's just a prop passing hell.*/}
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
                        </div>
                    </div>
                </div>
            </div>
        </>
    )
}

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