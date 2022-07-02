import '../styles/globals.css'
import type {AppProps} from 'next/app'
import {ThemeProvider} from 'next-themes'
import Head from "next/head";
import Link from "next/link";
import GitHubIcon from "@mui/icons-material/GitHub";
import LinkedInIcon from "@mui/icons-material/LinkedIn";
import SiteFooter from "../components/SiteFooter"

// TODO: Laita sivu toimimaan hyvin mobilen kanssa.
// TODO: Loyda tapa kayttaa uudelleen yhta handleChange funkiota kaikkiin asetuksiin.
function MyApp({Component, pageProps}: AppProps) {
    return (
        <ThemeProvider defaultTheme="night">
            <Head>
                <title>Varaapoyta</title>
                <meta name="viewport" content="viewport-fit=cover"/>
            </Head>
            <Component {...pageProps} />
            <SiteFooter/>
        </ThemeProvider>

    )
}

export default MyApp
