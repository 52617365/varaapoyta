import React from "react";
import SpiceLevel from "./SpiceLevel"
import Button from "../components/Button"

function Form() {
    const [state, setState] = React.useState({
        tulinen: false,
        mieto: false,
        kirpea: false,
    });

    const {tulinen, mieto, kirpea} = state;

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setState({
            ...state,
            [event.target.name]: event.target.checked,
        });
    };

    return (
        <div className={"flex flex-col justify-center items-center h-full inset-x-0 top-0"}>
            <div className={"text-center border-4 rounded-md"}>
                <p>Tulisuus</p>
                <SpiceLevel tulinen={tulinen} mieto={mieto} kirpea={kirpea} handleChange={handleChange}/>
                <p>Add something</p>
                {/* TODO: Add other options here, E.g other user preferences. */}
                <SpiceLevel tulinen={tulinen} mieto={mieto} kirpea={kirpea} handleChange={handleChange}/>
                {/* TODO: Make button go back to false state when fetching is done.*/}
                <Button/>
            </div>
        </div>
    )
}

export default Form;