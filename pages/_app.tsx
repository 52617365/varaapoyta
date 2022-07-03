import '../styles/globals.css'
import type {AppProps} from 'next/app'
import {ThemeProvider} from 'next-themes'
import Head from "next/head";
import SiteFooter from "../components/SiteFooter"

// TODO: Laita sivu toimimaan hyvin mobilen kanssa.
// TODO: Loyda tapa kayttaa uudelleen yhta handleChange funkiota kaikkiin asetuksiin.
// TODO: Asetukset sivulla, laita asetusruudut sulkeutumaan silloin, kun niita painetaan uudestaan.

function MyApp({Component, pageProps}: AppProps) {
    return (
        <ThemeProvider defaultTheme="night">
            <Head>
                <title>Varaapoyta</title>
                <meta name="viewport" content="width=device-width,height=device-height,initial-scale=1.0"/>
            </Head>
            <Component {...pageProps} />
            <SiteFooter/>
        </ThemeProvider>

    )
}

export default MyApp
