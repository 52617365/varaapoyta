import React from "react";

function Button({text}: {text: string}) {
    const [buttonLoading, setButtonLoading] = React.useState(false);

    function setButton() {
        setButtonLoading(true)
    }

    return (
        <>
            {buttonLoading ? <button className="btn loading">Ladataan...</button> :
                <button onClick={setButton} className="btn">{text}</button>}
        </>
    )
}
export default Button;