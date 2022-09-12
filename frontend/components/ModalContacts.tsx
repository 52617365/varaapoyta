import api_response from "../interfaces/api_response_interface";
import Countdown from "../components/Countdown";

function get_open_restaurant_times(ravintola: api_response) {
  if (ravintola.openingTime.time_till_restaurant_closed_hours == -1) {
    return <p>CLOSED!</p>;
  } else {
    return (
      <Countdown
        hours={ravintola.openingTime.time_till_restaurant_closed_hours}
        minutes={ravintola.openingTime.time_till_restaurant_closed_minutes}
        seconds={0}
      />
    );
  }
}

function get_open_restaurant_kitchen_times(ravintola: api_response) {
  if (ravintola.openingTime.time_left_to_reserve_hours == -1) {
    return <p>CLOSED!</p>;
  } else {
    return (
      <Countdown
        hours={ravintola.openingTime.time_left_to_reserve_hours}
        minutes={ravintola.openingTime.time_left_to_reserve_minutes}
        seconds={0}
      />
    );
  }
}

function ModalInformation({ ravintola }: { ravintola: api_response }) {
  return (
    <>
      <label
        htmlFor={"information" + ravintola.id}
        className="btn modal-button"
      >
        Tiedot
      </label>
      <input
        type="checkbox"
        id={"information" + ravintola.id}
        className="modal-toggle"
      />
      <label
        htmlFor={"information" + ravintola.id}
        className="modal cursor-pointer"
      >
        <label className="modal-box relative" htmlFor="">
          <p className="py-4">
            Kaupunki: {ravintola.address.municipality.fi_FI}
          </p>
          <p className="py-4">Osoite: {ravintola.address.street.fi_FI}</p>
          <p className="py-4">Postinumero: {ravintola.address.zipCode}</p>
          <p className="py-4">
            Aukioloajat: {ravintola.openingTime.restaurantTime.ranges[0].start}-
            {ravintola.openingTime.restaurantTime.ranges[0].end}
          </p>
          <p className="py-4">
            Ravintolan sulkeutuminen
            {get_open_restaurant_times(ravintola)}
          </p>
          <p className="py-4">
            Varaukseen jäljellä olevaa aikaa
            <div className="content-center">
              {get_open_restaurant_kitchen_times(ravintola)}
            </div>
          </p>
        </label>
      </label>
    </>
  );
}

export default ModalInformation;
