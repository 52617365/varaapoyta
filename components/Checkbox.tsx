import React from 'react';

function Checkbox({name}: {name: string}) {
    return (
        <div className={"form-control"}>
            <label className="label cursor-pointer">
                <span className="label-text">{name}</span>
                {/*defaultChecked={mieto} onChange={handleChange}*/}
                <input type="checkbox" name={name}
                       className="checkbox"/>
            </label>
        </div>
    );
}

export default Checkbox;