import SiteHead from "../components/SiteHead";
import SiteFooter from "../components/SiteFooter";
import React from "react";
import DropDown from "../components/Dropdown";
import Checkbox from "../components/Checkbox";
import SubmitButton from "../components/SubmitButton"

function GenerateCheckboxes(items: string[]) {
    return (
        items.map((item) => {
            return (
                <Checkbox key={item} name={item}/>
            )
        }))
}

function Settings({restaurants}: { restaurants: string[] }) {
    const [settings, setSettings] = React.useState({
        restaurants: [],
    })

    const restaurantBoxes = GenerateCheckboxes(restaurants)

    const handleChange = (e: any) => {
        // Destructuring
        const {name, checked} = e.target;
        const {restaurants} = settings;

        // console.log(`${name} is ${checked}`);

        // Case 1 : The user checks the box
        if (checked) {
            setSettings({
                // @ts-ignore
                restaurants: [...restaurants, name],
            });
        }
        // Case 2  : The user unchecks the box
        else {
            setSettings({
                restaurants: restaurants.filter((e) => e !== name),
            });
        }

    };

    console.log(settings.restaurants)

    return (
        <>
            <SiteHead/>
            <div className={"flex flex-col justify-center items-center h-full inset-x-0 top-0"}>
                {/*TODO: Capture state of selected items.*/}
                {/*Ravintolat settings*/}
                <form onChange={handleChange}>
                    <DropDown
                        items={restaurantBoxes}
                        name={"Ravintolat"}
                    />
                    <div className={"absolute bottom-1/3"}>
                        <SubmitButton text={"Tallenna asetukset"}/>
                    </div>
                </form>
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