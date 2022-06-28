import React from "react";

function Spicyness(props: any) {
    const {tulinen, mieto, kirpea, handleChange} = props;
    return (
        <div
            className="grid gap-3 grid-cols-3">
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
    )
}

export default Spicyness;