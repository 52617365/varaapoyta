import Link from "next/link";
import api_response from "../interfaces/api_response_interface";

function ModalLinks({ ravintola }: { ravintola: api_response }) {
    return (
        <>
            <label htmlFor={"linkit" + ravintola.restaurant.id} className="btn modal-button">Linkit</label>
            <input type="checkbox" id={"linkit" + ravintola.restaurant.id} className="modal-toggle" />
            <label htmlFor={"linkit" + ravintola.restaurant.id} className="modal cursor-pointer">
                <label className="modal-box relative" htmlFor="">
                    <p className="py-4">
                        <Link href={ravintola.restaurant.links.tableReservationLocalized.fi_FI}>
                            <a target="_blank" className="text-sky-600">Varaussivulle</a>
                        </Link>
                    </p>
                    <p className="py-4">
                        <Link href={ravintola.restaurant.links.homepageLocalized.fi_FI}>
                            <a target="_blank" className="text-sky-600">Kotisivulle</a>
                        </Link>
                    </p>
                </label>
            </label>
        </>
    );
}

export default ModalLinks
