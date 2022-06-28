import type {NextPage} from 'next'
import SiteFooter from "../components/SiteFooter"
import SiteHead from "../components/SiteHead";
import Form from "../components/Form";
import Introduction from "../components/Introduction"

const Home: NextPage = () => {
    return (
        <>
            <SiteHead/>
            <Introduction/>
            <Form/>
            <SiteFooter/>
        </>
    )
}
export default Home
