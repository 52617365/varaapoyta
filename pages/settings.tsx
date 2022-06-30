import SiteHead from "../components/SiteHead";
import SiteFooter from "../components/SiteFooter";
import React from "react";
import DropDown from "../components/Dropdown";
import Checkbox from "../components/Checkbox";

function Settings({restaurants}: { restaurants: string[] }) {
    const [valitutRavintolat, lisaaRavintola] = React.useState(
        []
    );
    const handleRavintola = (newValue: String) => {
        // @ts-ignore
        lisaaRavintola((ravintolat: String) => [...ravintolat, newValue]);
    };

    const restaurantBoxes = restaurants.map((restaurant) => {
        return (
            // eslint-disable-next-line react/jsx-key
            <Checkbox name={restaurant}/>
        )
    })
    return (
        <>
            <SiteHead/>
            <div className={"flex flex-col justify-center items-center h-full inset-x-0 top-0"}>
                {/*TODO: Capture state of selected items.*/}

                {/*Ravinolat settings*/}
                <DropDown
                    handleChange={"asd"}
                    items={restaurantBoxes}
                    name={"Ravintolat"}
                />
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