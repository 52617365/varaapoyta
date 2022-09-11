import api_response from "../interfaces/api_response_interface";
import Countdown from "../components/Countdown";

function ModalInformation({ ravintola}: { ravintola: api_response}) {
  return (
    <>
      <label htmlFor={"information" + ravintola.id} className="btn modal-button">
        Tiedot
      </label>
      <input type="checkbox" id={"information" + ravintola.id} className="modal-toggle" />
      <label htmlFor={"information" + ravintola.id} className="modal cursor-pointer">
        <label className="modal-box relative" htmlFor="">
          <p className="py-4">
            Kaupunki: {ravintola.address.municipality.fi_FI}
          </p>
          <p className="py-4">Osoite: {ravintola.address.street.fi_FI}</p>
          <p className="py-4">Postinumero: {ravintola.address.zipCode}</p>
          <p className="py-4">
            Aukioloajat:{" "}
            {ravintola.openingTime.restaurantTime.ranges[0].start}-
            {ravintola.openingTime.restaurantTime.ranges[0].end}
          </p>
          <p className="py-4">
            Ravintolan sulkeutuminen
            <Countdown
              hours={ravintola.openingTime.time_till_restaurant_closed_hours}
              minutes={
                ravintola.openingTime.time_till_restaurant_closed_minutes
              }
              seconds={0}
            />
          </p>
          <p className="py-4">
            Ravintolan keittiön sulkeutuminen
            <div className="content-center">
              <Countdown
                hours={ravintola.openingTime.time_till_kitchen_closed_hours}
                minutes={
                  ravintola.openingTime.time_till_kitchen_closed_minutes
                }
                seconds={0}
              />
            </div>
          </p>
        </label>
      </label>
    </>
  );
}

export default ModalInformation;
