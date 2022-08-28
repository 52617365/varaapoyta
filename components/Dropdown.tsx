import {FormEventHandler} from "react";

function DropDown({
                      name,
                      items,
                      onChange
                  }: { name: string, items: Array<JSX.Element>, onChange: FormEventHandler<HTMLDivElement> }) {
    return (
        <div onChange={onChange}
             className="dropdown dropdown-right dropdown-hover">
            <label tabIndex={0} className="btn m-1">{name}</label>
            <ul tabIndex={0} className="dropdown-content menu shadow bg-base-100 rounded-box">
                {items}
            </ul>
        </div>
    )
}

export default DropDown;
