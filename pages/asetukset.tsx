import SiteHead from "../components/SiteHead";
import SiteFooter from "../components/SiteFooter";
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

function Asetukset({restaurants}: { restaurants: string[] }) {
    const [settings, setSettings] = React.useState({
        restaurants: [],
        //    TODO: Insert other wanted settings in here and do the same for all of them.
    })

    const restaurantBoxes = GenerateCheckboxes(restaurants)

    const handleChange = (e: any) => {
        // Destructuring
        const {name, checked} = e.target;
        const {restaurants} = settings;

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
            <div className="relative hero min-h-screen bg-base-200">
                <div className="hero-content text-center">
                    <div className="absolute top-1/3 max-w-md pr-5">
                        <h1 className={"text-xl pb-10"}>Asetukset</h1>
                        <div onChange={handleChange} className={"grid gap-10"}>
                            <DropDown
                                items={restaurantBoxes}
                                name={"Ravintolat"}
                            />
                            {/*TODO: Once user has saved settings, make button go back into the false state (not loading)*/}
                            <Button text={"Tallenna asetukset"}/>
                            <Link href={"/varaa"}>
                                <a>
                                    <Button text={"Takaisin varaamaan"}/>
                                </a>
                            </Link>
                        </div>
                    </div>
                    <SiteFooter/>
                </div>
            </div>
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

export default Asetukset;