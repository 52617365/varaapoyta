import api_response from "../interfaces/api_response_interface";

function ModalOpenTimes({ ravintola }: { ravintola: api_response }) {
  return (
    <>
      <label
        htmlFor={"vapaat_ajat" + ravintola.id}
        className="btn modal-button"
      >
        Vapaat pöydät
      </label>
      <input
        type="checkbox"
        id={"vapaat_ajat" + ravintola.id}
        className="modal-toggle"
      />
      <label
        htmlFor={"vapaat_ajat" + ravintola.id}
        className="modal cursor-pointer"
      >
        <label className="modal-box relative" htmlFor="">
          <div className="grid grid-cols-4 gap-4">
            {ravintola.availableTables.map((available_time_slot: string) => {
              return (
                <div key={available_time_slot}>
                  <button className="btn btn-primary w-full">
                    {available_time_slot}
                  </button>
                </div>
              );
            })}
          </div>
        </label>
      </label>
    </>
  );
}

export default ModalOpenTimes;
