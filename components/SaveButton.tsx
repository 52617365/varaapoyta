import React, {MouseEventHandler} from "react";

function SaveButton({text, buttonLoading, setButton}: {text: string, buttonLoading: any, setButton: MouseEventHandler}) {
    return (
        <>
            {buttonLoading ? <button className="btn loading">Ladataan...</button> :
                <button onClick={setButton} className="btn">{text}</button>}
        </>
    )
}

export default SaveButton;