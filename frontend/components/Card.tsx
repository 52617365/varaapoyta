import api_response from "../interfaces/api_response_interface";
import ModalInformation from "../components/ModalContacts";
import ModalLinks from "../components/ModalLinks";
import ModalOpenTimes from "../components/ModalOpenTimes";

function Card({ ravintola }: { ravintola: api_response }) {
  return (
    <div className="card w-150 bg-base-100 shadow-xl">
      <div className="card-body items-center text-center">
        <h2 className="card-title">{ravintola.name.fi_FI}</h2>
        <div className="pr-4 space-x-2">
          <ModalInformation information={ravintola} />
          <ModalLinks information={ravintola} />
          <ModalOpenTimes information={ravintola} />
        </div>
        <div className="card-actions justify-end">
          <button
            className="btn btn-disabled"
            tabIndex={-1}
            role="button"
            aria-disabled="true"
          >
            Varaa (V2)
          </button>
        </div>
      </div>
    </div>
  );
}

export default Card;
