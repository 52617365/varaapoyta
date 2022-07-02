import React from "react";

function SubmitButton({text}: { text: string }) {
    const [buttonLoading, setButtonLoading] = React.useState(false);

    function setButton() {
        setButtonLoading(!buttonLoading)
    }

    return (
        <>
                {buttonLoading ? <button className="btn loading">Ladataan...</button> :
                    <button onClick={setButton} form="form" type="submit" className="btn">{text}</button>}
        </>
    )
}
export default SubmitButton;