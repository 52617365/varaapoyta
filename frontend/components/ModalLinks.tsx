import Link from "next/link";
import api_response from "../interfaces/api_response_interface";

function ModalLinks({ information }: { information: api_response }) {
    return (
        <>
            <label htmlFor="my-modal-5" className="btn modal-button">Linkit</label>
            <input type="checkbox" id="my-modal-5" className="modal-toggle" />
            <label htmlFor="my-modal-5" className="modal cursor-pointer">
                <label className="modal-box relative" htmlFor="">
                    <p className="py-4">
                        <Link href={information.links.tableReservationLocalized.fi_FI}>
                            <a target="_blank" className="text-sky-600">Varaussivulle</a>
                        </Link>
                    </p>
                    <p className="py-4">
                        <Link href={information.links.homepageLocalized.fi_FI}>
                            <a target="_blank" className="text-sky-600">Kotisivulle</a>
                        </Link>
                    </p>
                </label>
            </label>
        </>
    );
}

export default ModalLinks
