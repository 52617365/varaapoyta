import Link from "next/link";
import {ApiResponse} from "../interfaces/api_response_interface";

function ModalLinks({ restaurant }: { restaurant: ApiResponse }) {
    return (
        <>
            <label htmlFor={"linkit" + restaurant.restaurant.id} className="btn modal-button">Linkit</label>
            <input type="checkbox" id={"linkit" + restaurant.restaurant.id} className="modal-toggle" />
            <label htmlFor={"linkit" + restaurant.restaurant.id} className="modal cursor-pointer">
                <label className="modal-box relative" htmlFor="">
                    <p className="py-4">
                        <Link href={restaurant.restaurant.Links.tableReservationLocalized.fi_FI}>
                            <a target="_blank" className="text-sky-600">Varaussivulle</a>
                        </Link>
                    </p>
                    <p className="py-4">
                        <Link href={restaurant.restaurant.Links.homepageLocalized.fi_FI}>
                            <a target="_blank" className="text-sky-600">Kotisivulle</a>
                        </Link>
                    </p>
                </label>
            </label>
        </>
    );
}

export default ModalLinks
