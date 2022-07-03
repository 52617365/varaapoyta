import React from "react";
import DropDown from "../components/Dropdown";
import Checkbox from "../components/Checkbox";
import Button from "../components/Button";
import Link from "next/link";

function GenerateCheckboxes(items: string[]) {
    return (
        items.map((item) => {
            return (
                <Checkbox key={item} name={item}/>
            )
        }))
}

function DropDownHandler(e: any, key: any, settings: any, setSettings: any) {
    // Destructuring
    const {name, checked} = e.target;

    // Case 1 : The user checks the box
    if (checked) {
        setSettings((prev: any) => ({
            ...prev,
            [key]: [...settings[key], name]
        }))
        // Case 2  : The user unchecks the box
    } else {
        setSettings({
            [key]: settings[key].filter((e: string) => e !== name)
        });
    }

    // Case 1 : The user checks the box
    // if (checked) {
    //     setSettings({
    //         //@ts-ignore
    //         settings["name"]: [...items, name];
    //     })
    // }
    // // Case 2  : The user unchecks the box
    // else {
    //     setSettings({
    //         key: items.filter((e: string) => e !== name)
    //     });
    // }
}

function Asetukset({ravintolat, kaupungit}: { ravintolat: string[], kaupungit: string[] }) {
    const [settings, setSettings] = React.useState({
        ravintolat: [],
        kaupungit: [],
        //    TODO: Insert other wanted settings in here and do the same for all of them.
    })

    const ravintolaBoxes = GenerateCheckboxes(ravintolat)
    const kaupungitBoxes = GenerateCheckboxes(kaupungit)

    // TODO: Make this function generic so that it works with all settings.

    console.log(settings)

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
                                onChange={val => DropDownHandler(val, "ravintolat", settings, setSettings)}
                                items={ravintolaBoxes}
                                name={"Ravintolat"}
                            />
                            <DropDown
                                onChange={val => DropDownHandler(val, "kaupungit", settings, setSettings)}
                                items={kaupungitBoxes}
                                name={"Kaupungit"}
                            />
                            <Button text={"Tallenna asetukset"}/>
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
    // restaurants variableen tulee tulevaisuudessa tietokannasta tieto. se kyllakin staattisesti renderoityna at compile time.
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