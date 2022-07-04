import '../styles/globals.css'
import type {AppProps} from 'next/app'
import {ThemeProvider} from 'next-themes'
import Head from "next/head";
import SiteFooter from "../components/SiteFooter"

// TODO: Instead of logging in, use local storage to store settings.
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
