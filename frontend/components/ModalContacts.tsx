import api_response from "../interfaces/api_response_interface";
import Countdown from "../components/Countdown";

function ModalInformation({ information }: { information: api_response }) {
  return (
    <>
      <label htmlFor="my-modal-4" className="btn modal-button">
        Tiedot
      </label>
      <input type="checkbox" id="my-modal-4" className="modal-toggle" />
      <label htmlFor="my-modal-4" className="modal cursor-pointer">
        <label className="modal-box relative" htmlFor="">
          <p className="py-4">
            Kaupunki: {information.address.municipality.fi_FI}
          </p>
          <p className="py-4">Osoite: {information.address.street.fi_FI}</p>
          <p className="py-4">Postinumero: {information.address.zipCode}</p>
          <p className="py-4">
            Aukioloajat:{" "}
            {information.openingTime.restaurantTime.ranges[0].start}-
            {information.openingTime.restaurantTime.ranges[0].end}
          </p>
          <p className="py-4">
            Ravintolan sulkeutuminen
            <Countdown
              hours={information.openingTime.time_till_restaurant_closed_hours}
              minutes={
                information.openingTime.time_till_restaurant_closed_minutes
              }
              seconds={0}
            />
          </p>
          <p className="py-4">
            Ravintolan keittiön sulkeutuminen
            <div className="content-center">
              <Countdown
                hours={information.openingTime.time_till_kitchen_closed_hours}
                minutes={
                  information.openingTime.time_till_kitchen_closed_minutes
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
