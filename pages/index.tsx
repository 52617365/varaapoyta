import type {NextPage} from 'next'
import Link from "next/link"

const Home: NextPage = () => {
    return (
        <Link className="relative" href={"/login"}>
            <button
                className="absolute m-auto top-0 left-0 right-0 bottom-0 btn btn-xs sm:btn-sm md:btn-md lg:btn-lg">Kirjaudu
                sisään
                tallentaaksesi mieltymyksesi
            </button>
        </Link>
    )
}
export default Home
