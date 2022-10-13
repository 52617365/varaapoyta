import api_response from "../interfaces/api_response_interface";
import Countdown from "../components/Countdown";

function get_open_restaurant_times(apiResponse: api_response) {
  if (apiResponse.restaurant.openingTime.time_till_restaurant_closed_hours == -1) {
    return <p>CLOSED!</p>;
  } else {
    return (
      <Countdown
        hours={apiResponse.restaurant.openingTime.time_till_restaurant_closed_hours}
        minutes={apiResponse.restaurant.openingTime.time_till_restaurant_closed_minutes}
        seconds={0}
      />
    );
  }
}

function get_open_restaurant_kitchen_times(apiResponse: api_response) {
  if (apiResponse.restaurant.openingTime.time_left_to_reserve_hours == -1) {
    return <p>CLOSED!</p>;
  } else {
    return (
      <Countdown
        hours={apiResponse.restaurant.openingTime.time_left_to_reserve_hours}
        minutes={apiResponse.restaurant.openingTime.time_left_to_reserve_minutes}
        seconds={0}
      />
    );
  }
}

function ModalInformation({ apiResponse }: { apiResponse: api_response }) {
  return (
    <>
      <label
        htmlFor={"information" + apiResponse.restaurant.id}
        className="btn modal-button"
      >
        Tiedot
      </label>
      <input
        type="checkbox"
        id={"information" + apiResponse.restaurant.id}
        className="modal-toggle"
      />
      <label
        htmlFor={"information" + apiResponse.restaurant.id}
        className="modal cursor-pointer"
      >
        <label className="modal-box relative" htmlFor="">
          <p className="py-4">
            Kaupunki: {apiResponse.restaurant.address.municipality.fi_FI}
          </p>
          <p className="py-4">Osoite: {apiResponse.restaurant.address.street.fi_FI}</p>
          <p className="py-4">Postinumero: {apiResponse.restaurant.address.zipCode}</p>
          <p className="py-4">
            Aukioloajat: {apiResponse.restaurant.openingTime.restaurantTime.ranges[0].start}-
            {apiResponse.restaurant.openingTime.restaurantTime.ranges[0].end}
          </p>
          <p className="py-4">
            Ravintolan sulkeutuminen
            {get_open_restaurant_times(apiResponse)}
          </p>
          <p className="py-4">
            Varaukseen jäljellä olevaa aikaa
            <div className="content-center">
              {get_open_restaurant_kitchen_times(apiResponse)}
            </div>
          </p>
        </label>
      </label>
    </>
  );
}

export default ModalInformation;
