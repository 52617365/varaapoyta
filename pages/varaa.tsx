import Button from "../components/Button"
import SiteFooter from "../components/SiteFooter";
import SiteHead from "../components/SiteHead";

function Varaa() {
    return (
        <>
           <SiteHead/>
            <div className="hero min-h-screen bg-base-200">
                <div className="hero-content text-center">
                    <div className="max-w-md pr-5">
                        {/*TODO: Add different paths to different oauth logins here*/}
                        <h1 className={"text-xl pb-5"}>Raflaamo varaaja</h1>
                        <Button/>
                    </div>
                </div>
            </div>
            <SiteFooter/>
        </>
    )
}


export default Varaa;
