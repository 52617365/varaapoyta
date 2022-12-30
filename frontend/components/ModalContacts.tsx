import {ApiResponse} from "../interfaces/api_response_interface";
import Countdown from "../components/Countdown";

function get_open_restaurant_times(apiResponse: ApiResponse) {
  if (apiResponse.timeTillRestaurantCloses.hour == -1 || apiResponse.timeTillRestaurantCloses.minute == -1) {
    return <p>CLOSED!</p>;
  } else {
    return (
      <Countdown
        hours={apiResponse.timeTillRestaurantCloses.hour}
        minutes={apiResponse.timeTillRestaurantCloses.minute}
        seconds={0}
      />
    );
  }
}

function get_open_restaurant_kitchen_times(apiResponse: ApiResponse) {
  if (apiResponse.timeTillKitchenCloses.hour == -1 || apiResponse.timeTillKitchenCloses.minute == -1) {
    return <p>CLOSED!</p>;
  } else {
    return (
      <Countdown
        hours={apiResponse.timeTillKitchenCloses.hour}
        minutes={apiResponse.timeTillKitchenCloses.minute}
        seconds={0}
      />
    );
  }
}

function ModalInformation({ restaurant }: { restaurant: ApiResponse }) {
  return (
    <>
      <label
        htmlFor={"information" + restaurant.restaurant.id}
        className="btn modal-button"
      >
        Tiedot
      </label>
      <input
        type="checkbox"
        id={"information" + restaurant.restaurant.id}
        className="modal-toggle"
      />
      <label
        htmlFor={"information" + restaurant.restaurant.id}
        className="modal cursor-pointer"
      >
        <label className="modal-box relative" htmlFor="">
          <p className="py-4">
            Kaupunki: {restaurant.restaurant.address.municipality.fi_FI}
          </p>
          <p className="py-4">Osoite: {restaurant.restaurant.address.street.fi_FI}</p>
          <p className="py-4">Postinumero: {restaurant.restaurant.address.zipCode}</p>
          <p className="py-4">
            Aukioloajat: {restaurant.restaurant.openingTime.restaurantTime.ranges[0].start}-
            {restaurant.restaurant.openingTime.restaurantTime.ranges[0].end}
          </p>
          <p className="py-4">
            Ravintolan sulkeutuminen
            {get_open_restaurant_times(restaurant)}
          </p>
          <p className="py-4">
            Varaukseen jäljellä olevaa aikaa
            <div className="content-center">
              {get_open_restaurant_kitchen_times(restaurant)}
            </div>
          </p>
        </label>
      </label>
    </>
  );
}

export default ModalInformation;
