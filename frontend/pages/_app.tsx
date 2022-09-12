import "../styles/globals.css";
import type { AppProps } from "next/app";
import { ThemeProvider } from "next-themes";
import Head from "next/head";
import SiteFooter from "../components/SiteFooter";

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <>
      <ThemeProvider defaultTheme="night">
        <Head>
          <title>Varaapoyta V1</title>
          <meta
            name="viewport"
            content="width=device-width,height=device-height,initial-scale=1.0"
          />
        </Head>
        <Component {...pageProps} />
      </ThemeProvider>
      <ThemeProvider defaultTheme="night">
        <SiteFooter />
      </ThemeProvider>
    </>
  );
}

export default MyApp;
