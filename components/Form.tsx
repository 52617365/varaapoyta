import React from "react";

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
    const [buttonLoading, setButtonLoading] = React.useState(false);

    function setButton() {
        setButtonLoading(true)
    }

    return (
        <div className={"flex flex-col justify-center items-center h-full inset-x-0 top-0"}>
            <div className={"text-center border-4 rounded-md border-indigo-600 "}>
                <div
                    className="grid gap-3 grid-cols-3 grid-rows-3">
                    <div className={"form-control"}>
                        <label className="label cursor-pointer">
                            <span className="label-text">Tulinen</span>
                            <input type="checkbox" defaultChecked={tulinen} onChange={handleChange} name={"tulinen"}
                                   className="checkbox"/>
                        </label>
                    </div>
                    <div className={"form-control"}>
                        <label className="label cursor-pointer">
                            <span className="label-text">Mieto</span>
                            <input type="checkbox" defaultChecked={mieto} onChange={handleChange} name={"mieto"}
                                   className="checkbox"/>
                        </label>
                    </div>
                    <div className={"form-control"}>
                        <label className="label cursor-pointer">
                            <span className="label-text">Kirpeä</span>
                            <input type="checkbox" defaultChecked={kirpea} onChange={handleChange} name={"kirpea"}
                                   className="checkbox"/>
                        </label>
                    </div>
                </div>
                {/* TODO: Make button go back to false state when fetching is done.*/}
                {buttonLoading ? <button className="btn loading">Ladataan...</button> :
                    <button onClick={setButton} type="submit" className="btn">Etsi</button>}
            </div>
        </div>
    )
}

export default Form;