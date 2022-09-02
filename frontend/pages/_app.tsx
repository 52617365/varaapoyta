import '../styles/globals.css'
import type {AppProps} from 'next/app'
import {ThemeProvider} from 'next-themes'
import Head from "next/head";
import SiteFooter from "../components/SiteFooter"
import {SWRConfig} from 'swr'

function MyApp({Component, pageProps}: AppProps) {
    return (
        <SWRConfig value={{provider: () => new Map()}}>
            <ThemeProvider defaultTheme="night">
                <Head>
                    <title>Varaapoyta</title>
                    <meta name="viewport" content="width=device-width,height=device-height,initial-scale=1.0"/>
                </Head>
                <Component {...pageProps} />
                    <SiteFooter/>
            </ThemeProvider>
        </SWRConfig>
    )
}

export default MyApp
