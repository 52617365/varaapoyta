import FormControlLabel from '@mui/material/FormControlLabel';
import React from "react";
import Checkbox from "@mui/material/Checkbox";
import FavoriteBorder from "@mui/icons-material/FavoriteBorder";
import Favorite from "@mui/icons-material/Favorite";

function Form() {
    const [state, setState] = React.useState({
        tulinen: false,
        mieto: false,
        kirpea: false,
    });

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setState({
            ...state,
            [event.target.name]: event.target.checked,
        });
    };
    const {tulinen, mieto, kirpea} = state;

    return (
        <div className={"flex flex-col justify-center items-center"}>
            <div>
                <h1>Valitse mieltymyksiesi mukaan</h1>
                <div
                    className="text-center w-100 border-4 rounded-md border-indigo-600 grid gap-3 grid-cols-3 grid-rows-3">
                    <FormControlLabel
                        control={
                            <Checkbox icon={<FavoriteBorder/>} checkedIcon={<Favorite/>} checked={tulinen}
                                      onChange={handleChange} name="tulinen"/>
                        }
                        label="Tulinen"
                    />
                    <FormControlLabel
                        control={
                            <Checkbox icon={<FavoriteBorder/>} checkedIcon={<Favorite/>} checked={mieto}
                                      onChange={handleChange} name="mieto"/>
                        }
                        label="Mieto"
                    />
                    <FormControlLabel
                        control={
                            <Checkbox icon={<FavoriteBorder/>} checkedIcon={<Favorite/>} checked={kirpea}
                                      onChange={handleChange} name="kirpea"/>
                        }
                        label="Kirpeä"
                    />
                </div>
            </div>
        </div>
    )
}

export default Form;