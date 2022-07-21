import React, {DependencyList, EffectCallback, useEffect, useRef} from "react";
import DropDown from "../components/Dropdown";
import Checkbox from "../components/Checkbox";
import Button from "../components/Button";
import Link from "next/link";

// TODO: Laita local storage toimimaan.

// Locate checked checkboxes.
function GenerateCheckboxes(items: string[], checked_items: string[]) {
    return items.map((item) => {
        if (checked_items.includes(item)) {
            return <Checkbox name={item} checked={true}/>
        } else {
            return <Checkbox name={item} checked={false}/>
        }
    })
}

function RavintolaHandler(e: any, settings: any, setSettings: any) {
    // Destructuring
    const {name, checked} = e.target;

    // Case 1: The user checks the box
    if (checked) {
        // checking if the state alrdy has the value we're about to set to avoid duplicates.
        setSettings(
            [...settings, name],
        );
    }
    // Case 2: The user unchecks the box
    else {
        setSettings(
            settings.filter((e: string) => e !== name),
        );
    }
}

function KaupunkiHandler(e: any, settings: any, setSettings: any) {
    // Destructuring
    const {name, checked} = e.target;

    // Case 1: The user checks the box
    if (checked) {
        // checking if the state alrdy has the value we're about to set to avoid duplicates.
        setSettings(
            [...settings, name],
        );
    }
    // Case 2: The user unchecks the box
    else {
        setSettings(
            settings.filter((e: string) => e !== name),
        );
    }
}

/**
 * @param effect
 * @param dependencies
 *
 */
function useNoInitialEffect(
    effect: EffectCallback,
    dependencies?: DependencyList
) {
    //Preserving the true by default as initial render cycle
    const initialRender = useRef(true);

    useEffect(() => {
        let effectReturns: void | (() => void) = () => {
        };

        // Updating the ref to false on the first render, causing
        // subsequent render to execute the effect
        if (initialRender.current) {
            initialRender.current = false;
        } else {
            effectReturns = effect();
        }

        // Preserving and allowing the Destructor returned by the effect
        // to execute on component unmount and perform cleanup if
        // required.
        if (effectReturns && typeof effectReturns === 'function') {
            return effectReturns;
        }
        return undefined;
    }, [effect]);
}

function Asetukset({ravintolat, kaupungit}: { ravintolat: string[], kaupungit: string[] }) {
    const [ravintola_lista, lisaaRavintola] = React.useState([])
    const [kaupunki_lista, lisaaKaupunki] = React.useState([])

    let ravintolaBoxes = GenerateCheckboxes(ravintolat, ravintola_lista)
    let kaupungitBoxes = GenerateCheckboxes(kaupungit, kaupunki_lista)

    console.log(ravintola_lista)
    console.log(kaupunki_lista)
    // loading up states on page load.
    useEffect(() => {
        let ravintolat_storage = localStorage.getItem('varaapoyta_ravintolat')
        let kaupungit_storage = localStorage.getItem('varaapoyta_kaupungit')

        if (ravintolat_storage != null) {
            let ravintolat_parsed = JSON.parse(ravintolat_storage)
            lisaaRavintola(ravintolat_parsed);
        }
        if (kaupungit_storage != null) {
            let kaupungit_parsed = JSON.parse(kaupungit_storage)
            lisaaKaupunki(kaupungit_parsed);
        }
    }, []);

    // Loading up the states into local storage on change.
    useNoInitialEffect(() => {
        // to avoid setting empty array into local storage on page load.
        // TODO: this is buggy cuz now it wont let me update local storage into 0 options.
        localStorage.setItem("varaapoyta_ravintolat", JSON.stringify(ravintola_lista))
        localStorage.setItem("varaapoyta_kaupungit", JSON.stringify(kaupunki_lista))
    }, [ravintola_lista, kaupunki_lista]);

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