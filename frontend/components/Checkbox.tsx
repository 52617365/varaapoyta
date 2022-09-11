import React from 'react';

function Checkbox({name, checked}: { name: string, checked: boolean }) {
    return (
        <div className={"form-control"}>
            <label className="label cursor-pointer">
                <span className="label-text">{name}</span>
                <input type="checkbox" defaultChecked={checked} name={name}
                       className="checkbox"/>
            </label>
        </div>
    );
}

export default Checkbox;