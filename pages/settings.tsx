import SiteHead from "../components/SiteHead";
import SiteFooter from "../components/SiteFooter";
import React from "react";
import DropDown from "../components/Dropdown";
import Checkbox from "../components/Checkbox";
import Button from "../components/Button";

function GenerateCheckboxes(items: string[], handler: Function) {
    return (
        items.map((item) => {
            return (
                <Checkbox key={item} handler={handler} name={item}/>
            )
        }))
}

function Settings({restaurants}: { restaurants: string[] }) {
    // in the end valitutRavintolat will contain all the user selected restaurants.
    // valitutAsetukset will be an object of arrays, array will contain all of the settings from the corresponding key (key will be a setting section)
    // E.g {restaurants: [res1, res2...]},
    const [valitutAsetukset, lisaaAsetus] = React.useState(
        {}
    );

    function asetusHandler(key: string, asetus: string) {
    //    TODO: Add functionality for settings
    }

    const restaurantBoxes = GenerateCheckboxes(restaurants, asetusHandler)

    return (
        <>
            <SiteHead/>
            <div className={"flex flex-col justify-center items-center h-full inset-x-0 top-0"}>
                {/*TODO: Capture state of selected items.*/}
                {/*Ravintolat settings*/}
                <DropDown
                    items={restaurantBoxes}
                    name={"Ravintolat"}
                />
                <div className={"absolute bottom-1/3"}>
                    <Button text={"Tallenna asetukset"}/>
                </div>
            </div>
            <SiteFooter/>
        </>
    )
}

// TODO: Add other data in here too.
export async function getStaticProps() {
    // TODO: Hae ravintola nimet jostain ja anna ne main componenttiin tasta.
    const restaurants = ["restaurant1", "restaurant2", "restaurant3"];
    return {
        props: {
            restaurants
        },
    }
}

export default Settings;