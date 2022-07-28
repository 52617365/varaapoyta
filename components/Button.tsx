import React from "react";

// TODO: add actual functionality to this function. E.g fetching.
function Button({text, setButton, buttonLoading}: { text: string, setButton: any, buttonLoading: any }) {
    return (
        <>
            {buttonLoading ? <button className="btn loading">Ladataan...</button> :
                <button onClick={setButton} className="btn">{text}</button>}
        </>
    )
}

export default Button;
