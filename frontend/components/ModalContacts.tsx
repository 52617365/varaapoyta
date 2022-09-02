import api_response from "../interfaces/api_response_interface"

function ModalContacts({information}: { information: api_response }) {
    return (
        <>
            <label htmlFor="my-modal-4" className="btn modal-button">Tiedot</label>
            <input type="checkbox" id="my-modal-4" className="modal-toggle"/>
            <label htmlFor="my-modal-4" className="modal cursor-pointer">
                <label className="modal-box relative" htmlFor="">
                    <p className="py-4">Kaupunki: {information.address.municipality.fi_FI}</p>
                    <p className="py-4">Osoite: {information.address.street.fi_FI}</p>
                    <p className="py-4">Postinumero: {information.address.zipCode}</p>
                </label>
            </label>
        </>
    )
}

export default ModalContacts