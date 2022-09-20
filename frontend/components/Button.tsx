import React from "react";

function Button({
  text,
  setButton,
  buttonLoading,
  textfield_text,
}: {
  text: string;
  setButton: any;
  buttonLoading: any;
  textfield_text: string;
}) {
  return (
    <>
      {buttonLoading ? (
        <button className="btn loading">Ladataan...</button>
      ) : (
        <button onClick={() => setButton(textfield_text)} className="btn">
          {text}
        </button>
      )}
    </>
  );
}

export default Button;
