import React from "react";

function Button() {
    const [buttonLoading, setButtonLoading] = React.useState(false);

    function setButton() {
        setButtonLoading(true)
    }

    return (
        <>
            {buttonLoading ? <button className="btn loading">Ladataan...</button> :
                <button onClick={setButton} type="submit" className="btn">Etsi</button>}
        </>
    )
}
export default Button;