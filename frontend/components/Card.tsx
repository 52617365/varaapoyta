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
          {/* <button
            className="btn btn-disabled"
            tabIndex={-1}
            role="button"
            aria-disabled="true"
          >
            Varaa (V2)
          </button> */}
          {/* <Countdown
            hours={texts.openingTime.time_till_restaurant_closed_hours}
            minutes={texts.openingTime.time_till_restaurant_closed_minutes}
            seconds={0}
          /> */}
        </div>
      </div>
    </div>
  );
}

export default Card;
