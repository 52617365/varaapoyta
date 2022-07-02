function DropDown({name, items}: { name: string, items: Array<JSX.Element> }) {
    return (
        <div className="dropdown dropdown-right dropdown-hover">
            <label tabIndex={0} className="btn m-1">{name}</label>
            <ul tabIndex={0} className="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52">
                {/*TODO: Add all restaurants here*/}
                {items}
            </ul>
        </div>
    )
}

export default DropDown;
