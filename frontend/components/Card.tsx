import api_response from "../interfaces/api_response_interface"
import ModalContacts from "../components/ModalContacts"
import ModalLinks from "../components/ModalLinks"

function Card({texts}: { texts: api_response}) {
    return (
        <div className="card w-150 bg-base-100 shadow-xl">
            <div className="card-body items-center text-center">
                <h2 className="card-title">{texts.name.fi_FI}</h2>
                    <div className="pr-4 space-x-2">
                        <ModalContacts information={texts}/>
                        <ModalLinks information={texts}/>
                    </div>
            <div className="card-actions justify-end">
                {/*TODO: add some timer here or something?*/}
                {/*TODO: Add the possibility to reserve here.*/}
                <button className="btn btn-primary">Varaa</button>
            </div>
            </div>
        </div>
       )
}

export default Card