import api_response from "../interfaces/api_response_interface";
import ModalContacts from "../components/ModalContacts";
import ModalLinks from "../components/ModalLinks";
import ModalOpenTimes from "../components/ModalOpenTimes";

function Card({ texts }: { texts: api_response }) {
  return (
    <div className="card w-150 bg-base-100 shadow-xl">
      <div className="card-body items-center text-center">
        <h2 className="card-title">{texts.name.fi_FI}</h2>
        <div className="pr-4 space-x-2">
          <ModalContacts information={texts} />
          <ModalLinks information={texts} />
          <ModalOpenTimes information={texts} />
        </div>
        <div className="card-actions justify-end">
          {/* TODO: add some timer here or something?*/}
          {/* TODO: Add the possibility to reserve here.*/}
          <button
            className="btn btn-disabled"
            tabIndex={-1}
            role="button"
            aria-disabled="true"
          >
            Varaa
          </button>
          {/* <button className="btn block" disabled>
              Varaa
              <span className="text-sm font-serif">
                <br></br>Tulossa pian
              </span>
            </button> */}
        </div>
        {/* <button className="btn btn-primary">Varaa</button> */}
      </div>
    </div>
  );
}

export default Card;
