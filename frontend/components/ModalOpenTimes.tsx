import api_response from "../interfaces/api_response_interface"

function ModalOpenTimes({information}: { information: api_response }) {
    return (
        <>
            <label htmlFor="my-modal-6" className="btn modal-button">Vapaat pöydät</label>
            <input type="checkbox" id="my-modal-6" className="modal-toggle"/>
            <label htmlFor="my-modal-6" className="modal cursor-pointer">
                <label className="modal-box relative" htmlFor="">
                <div className="grid grid-cols-4 gap-4">
                    {information.available_time_slots.map((available_time_slot: string) => {
                      return (
                               <div key={available_time_slot}>
                                     <button className="btn btn-primary w-full">{available_time_slot}</button>
                               </div>
                            //         <p className="py-6">{available_time_slot}</p>
                            //         <button className="btn btn-primary">Varaa</button>
                            // </div>
                       )
                    })}
                </div>
                </label>
            </label>
        </>
    )
}

export default ModalOpenTimes