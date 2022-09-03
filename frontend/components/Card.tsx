import api_response from "../interfaces/api_response_interface";
import ModalInformation from "../components/ModalContacts";
import ModalLinks from "../components/ModalLinks";
import ModalOpenTimes from "../components/ModalOpenTimes";
import Countdown from "../components/Countdown";

function Card({ texts }: { texts: api_response }) {
  return (
    <div className="card w-150 bg-base-100 shadow-xl">
      <div className="card-body items-center text-center">
        <h2 className="card-title">{texts.name.fi_FI}</h2>
        <div className="pr-4 space-x-2">
          <ModalInformation information={texts} />
          <ModalLinks information={texts} />
          <ModalOpenTimes information={texts} />
        </div>
        <div className="card-actions justify-end">
          {/* TODO: add some timer here or something?*/}
          {/* TODO: Add the possibility to reserve here.*/}
          {/* <button className="btn btn-primary">Varaa</button> */}
          <button
            className="btn btn-disabled"
            tabIndex={-1}
            role="button"
            aria-disabled="true"
          >
            Varaa (V2)
          </button>
          {/* TODO: Add something related to closing time here */}
          {/* Like a countdown timer or something */}
          <Countdown hours={0} minutes={1} seconds={1} />
        </div>
      </div>
    </div>
  );
}

function get_time_remaining() {}
export default Card;
