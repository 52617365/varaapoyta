function DropDown() {
    return (
        <div className="dropdown">
            <label tabIndex={0} className="btn m-1">Sijanti</label>
            <ul tabIndex={0} className="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52">
                {/*TODO: Add all finnish locations here*/}
                <li><a>Item 1</a></li>
                <li><a>Item 2</a></li>
            </ul>
        </div>
    )
}
export default DropDown;