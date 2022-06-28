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
        <div className={"w-full"}>
            <div className="w-1/2 grid gap-4 grid-cols-3 grid-rows-3">
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
    )
}

export default Form;